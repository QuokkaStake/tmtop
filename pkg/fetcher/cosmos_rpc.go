package fetcher

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	configPkg "main/pkg/config"
	"main/pkg/http"
	"main/pkg/types"
	"net/url"
	"strconv"
	"strings"
	"sync"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/std"
	"github.com/rs/zerolog"

	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	queryTypes "github.com/cosmos/cosmos-sdk/types/query"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	providerTypes "github.com/cosmos/interchain-security/v3/x/ccv/provider/types"
)

type CosmosRPCDataFetcher struct {
	Config         configPkg.Config
	Logger         zerolog.Logger
	Client         *http.Client
	ProviderClient *http.Client

	Registry   codecTypes.InterfaceRegistry
	ParseCodec *codec.ProtoCodec
}

func NewCosmosRPCDataFetcher(config configPkg.Config, logger zerolog.Logger) *CosmosRPCDataFetcher {
	interfaceRegistry := codecTypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	parseCodec := codec.NewProtoCodec(interfaceRegistry)

	return &CosmosRPCDataFetcher{
		Config:         config,
		Logger:         logger.With().Str("component", "cosmos_data_fetcher").Logger(),
		ProviderClient: http.NewClient(logger, "cosmos_data_fetcher", config.ProviderRPCHost),
		Client:         http.NewClient(logger, "cosmos_data_fetcher", config.RPCHost),
		Registry:       interfaceRegistry,
		ParseCodec:     parseCodec,
	}
}

func (f *CosmosRPCDataFetcher) GetProviderOrConsumerClient() *http.Client {
	if f.Config.ProviderRPCHost != "" {
		return f.ProviderClient
	}

	return f.Client
}

func (f *CosmosRPCDataFetcher) GetValidatorAssignedConsumerKey(
	providerValcons string,
) (*providerTypes.QueryValidatorConsumerAddrResponse, error) {
	query := providerTypes.QueryValidatorConsumerAddrRequest{
		ChainId:         f.Config.ConsumerChainID,
		ProviderAddress: providerValcons,
	}

	var response providerTypes.QueryValidatorConsumerAddrResponse
	if err := f.AbciQuery(
		"/interchain_security.ccv.provider.v1.Query/QueryValidatorConsumerAddr",
		&query,
		&response,
		f.ProviderClient,
	); err != nil {
		return nil, err
	}

	return &response, nil
}

func (f *CosmosRPCDataFetcher) AbciQuery(
	method string,
	message codec.ProtoMarshaler,
	output codec.ProtoMarshaler,
	client *http.Client,
) error {
	dataBytes, err := message.Marshal()
	if err != nil {
		return err
	}

	methodName := fmt.Sprintf("\"%s\"", method)
	queryURL := fmt.Sprintf(
		"/abci_query?path=%s&data=0x%x",
		url.QueryEscape(methodName),
		dataBytes,
	)

	var response types.AbciQueryResponse
	if err := client.Get(queryURL, &response); err != nil {
		return err
	}

	if response.Result.Response.Code != 0 {
		return fmt.Errorf(
			"error in Tendermint response: expected code 0, but got %d, error: %s",
			response.Result.Response.Code,
			response.Result.Response.Log,
		)
	}

	return output.Unmarshal(response.Result.Response.Value)
}

func (f *CosmosRPCDataFetcher) ParseValidator(validator stakingTypes.Validator) (types.ChainValidator, error) {
	if err := validator.UnpackInterfaces(f.ParseCodec); err != nil {
		return types.ChainValidator{}, err
	}

	addr, err := validator.GetConsAddr()
	if err != nil {
		return types.ChainValidator{}, err
	}

	return types.ChainValidator{
		Moniker:    validator.GetMoniker(),
		Address:    fmt.Sprintf("%X", addr),
		RawAddress: addr.String(),
	}, nil
}

func (f *CosmosRPCDataFetcher) GetValidators() (*types.ChainValidators, error) {
	query := stakingTypes.QueryValidatorsRequest{
		Pagination: &queryTypes.PageRequest{
			Limit: 1000,
		},
	}

	var validatorsResponse stakingTypes.QueryValidatorsResponse
	if err := f.AbciQuery(
		"/cosmos.staking.v1beta1.Query/Validators",
		&query,
		&validatorsResponse,
		f.GetProviderOrConsumerClient(),
	); err != nil {
		if strings.Contains(err.Error(), " please wait for first block") {
			return f.GetGenesisValidators()
		}
		return nil, err
	}

	validators := make(types.ChainValidators, len(validatorsResponse.Validators))

	for index, validator := range validatorsResponse.Validators {
		if chainValidator, err := f.ParseValidator(validator); err != nil {
			return nil, err
		} else {
			validators[index] = chainValidator
		}
	}

	if !f.Config.IsConsumer() {
		return &validators, nil
	}

	// fetching assigned keys
	var wg sync.WaitGroup
	var mutex sync.Mutex

	for index, validator := range validators {
		wg.Add(1)
		go func(validator types.ChainValidator, index int) {
			defer wg.Done()
			assignedKey, err := f.GetValidatorAssignedConsumerKey(validator.RawAddress)

			if err != nil {
				f.Logger.Error().Err(err).Msg("Could not fetch assigned key")
				return
			}

			assignedKeyAsString := assignedKey.GetConsumerAddress()
			if assignedKeyAsString != "" {
				addr, _ := sdkTypes.ConsAddressFromBech32(assignedKeyAsString)

				mutex.Lock()
				validators[index].AssignedAddress = addr.String()
				validators[index].RawAssignedAddress = fmt.Sprintf("%X", addr)
				mutex.Unlock()
			}
		}(validator, index)
	}

	wg.Wait()

	return &validators, nil
}

func (f *CosmosRPCDataFetcher) GetGenesisValidators() (*types.ChainValidators, error) {
	f.Logger.Info().Msg("Fetching genesis validators...")

	genesisChunks := make([][]byte, 0)
	var chunk int64 = 0

	for {
		f.Logger.Info().Int64("chunk", chunk).Msg("Fetching genesis chunk...")
		genesisChunk, total, err := f.GetGenesisChunk(chunk)
		f.Logger.Info().Int64("chunk", chunk).Int64("total", total).Msg("Fetched genesis chunk...")
		if err != nil {
			return nil, err
		}

		genesisChunks = append(genesisChunks, genesisChunk)

		if chunk >= total-1 {
			break
		}

		chunk++
	}

	genesisBytes := bytes.Join(genesisChunks, []byte{})
	f.Logger.Info().Int("length", len(genesisBytes)).Msg("Fetched genesis")

	var genesisStruct types.Genesis

	if err := json.Unmarshal(genesisBytes, &genesisStruct); err != nil {
		f.Logger.Error().Err(err).Msg("Error unmarshalling genesis")
		return nil, err
	}

	var stakingGenesisState stakingTypes.GenesisState
	if err := f.ParseCodec.UnmarshalJSON(genesisStruct.AppState.Staking, &stakingGenesisState); err != nil {
		f.Logger.Error().Err(err).Msg("Error unmarshalling staking genesis state")
		return nil, err
	}

	f.Logger.Info().Int("validators", len(stakingGenesisState.Validators)).Msg("Genesis unmarshalled")

	// 1. Trying to fetch validators from staking module. Works for chain which did not start
	// from the first block but had their genesis as an export from older chain.
	if len(stakingGenesisState.Validators) > 0 {
		validators := make(types.ChainValidators, len(stakingGenesisState.Validators))
		for index, validator := range stakingGenesisState.Validators {
			if chainValidator, err := f.ParseValidator(validator); err != nil {
				return nil, err
			} else {
				validators[index] = chainValidator
			}
		}

		return &validators, nil
	}

	return nil, fmt.Errorf("genesis validators fetching is not yet supported")
}

func (f *CosmosRPCDataFetcher) GetGenesisChunk(chunk int64) ([]byte, int64, error) {
	var response types.TendermintGenesisChunkResponse
	if err := f.Client.Get(
		fmt.Sprintf("/genesis_chunked?chunk=%d", chunk),
		&response,
	); err != nil {
		return nil, 0, err
	}

	if response.Result == nil {
		return nil, 0, fmt.Errorf("malformed response from node")
	}

	total, err := strconv.ParseInt(response.Result.Total, 10, 64)
	if err != nil {
		return nil, 0, err
	}

	return response.Result.Data, total, nil
}

func (f *CosmosRPCDataFetcher) GetUpgradePlan() (*types.Upgrade, error) {
	query := upgradeTypes.QueryCurrentPlanRequest{}

	var response upgradeTypes.QueryCurrentPlanResponse
	if err := f.AbciQuery(
		"/cosmos.upgrade.v1beta1.Query/CurrentPlan",
		&query,
		&response,
		f.Client,
	); err != nil {
		return nil, err
	}

	if response.Plan == nil {
		return nil, errors.New("upgrade plan is not present")
	}

	return &types.Upgrade{
		Name:   response.Plan.Name,
		Height: response.Plan.Height,
	}, nil
}
