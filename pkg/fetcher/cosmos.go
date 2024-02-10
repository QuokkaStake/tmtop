package fetcher

import (
	"errors"
	"fmt"
	configPkg "main/pkg/config"
	"main/pkg/http"
	"main/pkg/types"
	"net/url"
	"sync"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/std"
	"github.com/rs/zerolog"

	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	queryTypes "github.com/cosmos/cosmos-sdk/types/query"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	providerTypes "github.com/cosmos/interchain-security/x/ccv/provider/types"
)

type CosmosDataFetcher struct {
	Config         configPkg.Config
	Logger         zerolog.Logger
	Client         *http.Client
	ProviderClient *http.Client

	Registry   codecTypes.InterfaceRegistry
	ParseCodec *codec.ProtoCodec
}

func NewCosmosDataFetcher(config configPkg.Config, logger zerolog.Logger) *CosmosDataFetcher {
	interfaceRegistry := codecTypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	parseCodec := codec.NewProtoCodec(interfaceRegistry)

	return &CosmosDataFetcher{
		Config:         config,
		Logger:         logger.With().Str("component", "cosmos_data_fetcher").Logger(),
		ProviderClient: http.NewClient(logger, "cosmos_data_fetcher", config.ProviderRPCHost),
		Client:         http.NewClient(logger, "cosmos_data_fetcher", config.RPCHost),
		Registry:       interfaceRegistry,
		ParseCodec:     parseCodec,
	}
}

func (f *CosmosDataFetcher) GetProviderOrConsumerClient() *http.Client {
	if f.Config.ProviderRPCHost != "" {
		return f.ProviderClient
	}

	return f.Client
}

func (f *CosmosDataFetcher) GetValidatorAssignedConsumerKey(
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

func (f *CosmosDataFetcher) AbciQuery(
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

func (f *CosmosDataFetcher) GetValidators() (*types.ChainValidators, error) {
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
		return nil, err
	}

	validators := make(types.ChainValidators, len(validatorsResponse.Validators))

	for index, validator := range validatorsResponse.Validators {
		if err := validator.UnpackInterfaces(f.ParseCodec); err != nil {
			return nil, err
		}

		addr, err := validator.GetConsAddr()
		if err != nil {
			return nil, err
		}

		validators[index] = types.ChainValidator{
			Moniker:    validator.GetMoniker(),
			Address:    fmt.Sprintf("%X", addr),
			RawAddress: addr.String(),
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

func (f *CosmosDataFetcher) GetUpgradePlan() (*types.Upgrade, error) {
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
