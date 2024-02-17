package fetcher

import (
	"errors"
	"fmt"
	configPkg "main/pkg/config"
	"main/pkg/http"
	"main/pkg/types"

	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/rs/zerolog"
)

type CosmosLcdDataFetcher struct {
	Config configPkg.Config
	Logger zerolog.Logger
	Client *http.Client

	Registry    codecTypes.InterfaceRegistry
	ParseCodec  *codec.ProtoCodec
	LegacyAmino *codec.LegacyAmino
}

func NewCosmosLcdDataFetcher(config configPkg.Config, logger zerolog.Logger) *CosmosLcdDataFetcher {
	interfaceRegistry := codecTypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	parseCodec := codec.NewProtoCodec(interfaceRegistry)

	return &CosmosLcdDataFetcher{
		Config:     config,
		Logger:     logger.With().Str("component", "cosmos_lcd_data_fetcher").Logger(),
		Client:     http.NewClient(logger, "cosmos_lcd_data_fetcher", config.LCDHost),
		Registry:   interfaceRegistry,
		ParseCodec: parseCodec,
	}
}

func (f *CosmosLcdDataFetcher) GetValidators() (*types.ChainValidators, error) {
	bytes, err := f.Client.GetPlain(
		"/cosmos/staking/v1beta1/validators?status=BOND_STATUS_BONDED&pagination.limit=1000",
	)

	if err != nil {
		return nil, err
	}

	var validatorsResponse stakingTypes.QueryValidatorsResponse

	if err := f.ParseCodec.UnmarshalJSON(bytes, &validatorsResponse); err != nil {
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
			Moniker:    validator.Description.Moniker,
			Address:    fmt.Sprintf("%X", addr),
			RawAddress: addr.String(),
		}
	}

	return &validators, nil
}

func (f *CosmosLcdDataFetcher) GetUpgradePlan() (*types.Upgrade, error) {
	var response upgradeTypes.QueryCurrentPlanResponse
	if err := f.Client.Get(
		"/cosmos/upgrade/v1beta1/current_plan",
		&response,
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
