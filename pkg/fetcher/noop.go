package fetcher

import "main/pkg/types"

type NoopDataFetcher struct {
}

func NewNoopDataFetcher() *NoopDataFetcher {
	return &NoopDataFetcher{}
}

func (f *NoopDataFetcher) GetValidators() (*types.ChainValidators, error) {
	return &types.ChainValidators{}, nil
}

func (f *NoopDataFetcher) GetUpgradePlan() (*types.Upgrade, error) {
	return nil, nil
}
