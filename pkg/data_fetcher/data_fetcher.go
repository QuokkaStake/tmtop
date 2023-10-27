package data_fetcher

import "main/pkg/types"

type DataFetcher interface {
	GetValidators() (types.ChainValidators, error)
}
