package fetcher

import (
	configPkg "main/pkg/config"
	"main/pkg/types"

	"github.com/rs/zerolog"
)

type DataFetcher interface {
	GetValidators() (*types.ChainValidators, error)
}

func GetDataFetcher(config configPkg.Config, logger zerolog.Logger) DataFetcher {
	if config.QueryValidators {
		return NewCosmosDataFetcher(config, logger)
	}

	return NewNoopDataFetcher()
}
