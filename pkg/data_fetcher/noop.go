package data_fetcher

import "main/pkg/types"

type NoopDataFetcher struct {
}

func NewNoopDataFetcher() *NoopDataFetcher {
	return &NoopDataFetcher{}
}

func (f *NoopDataFetcher) GetValidators() (types.ChainValidators, error) {
	return types.ChainValidators{}, nil
}
