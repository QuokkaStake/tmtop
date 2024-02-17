package fetcher

import (
	configPkg "main/pkg/config"
	"main/pkg/types"

	"github.com/rs/zerolog"
)

type DataFetcher interface {
	GetValidators() (*types.ChainValidators, error)
	GetUpgradePlan() (*types.Upgrade, error)
}

func GetDataFetcher(config configPkg.Config, logger zerolog.Logger) DataFetcher {
	if config.ChainType == "cosmos-rpc" {
		return NewCosmosRPCDataFetcher(config, logger)
	}

	if config.ChainType == "cosmos-lcd" {
		return NewCosmosLcdDataFetcher(config, logger)
	}

	return NewNoopDataFetcher()
}
