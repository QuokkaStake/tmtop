package data_fetcher

import (
	"github.com/rs/zerolog"
	configPkg "main/pkg/config"
	"main/pkg/types"
)

type DataFetcher interface {
	GetValidators() (types.ChainValidators, error)
}

func GetDataFetcher(config configPkg.Config, logger zerolog.Logger) DataFetcher {
	if config.QueryValidators {
		return NewCosmosDataFetcher(config, logger)
	}

	return NewNoopDataFetcher()
}
