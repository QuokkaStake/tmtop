package fetcher

import (
	configPkg "main/pkg/config"
	"main/pkg/types"

	"github.com/rs/zerolog"
)

type DataFetcher interface {
	GetValidators() (*types.ChainValidators, error)
	GetUpgradePlan() (*types.Upgrade, error)
	GetNetInfo(host string) (*types.NetInfo, error)
}

func GetDataFetcher(config *configPkg.Config, state *types.State, logger zerolog.Logger) DataFetcher {
	if config.ChainType == "tendermint" {
		return NewNoopDataFetcher()
	}

	if config.ChainType == "cosmos-lcd" {
		return NewCosmosLcdDataFetcher(config, logger)
	}

	return NewCosmosRPCDataFetcher(config, state, logger)
}
