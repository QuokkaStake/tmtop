package config

import "time"

type Config struct {
	RPCHost               string
	ProviderRPCHost       string
	ConsumerChainID       string
	RefreshRate           time.Duration
	ValidatorsRefreshRate time.Duration
	ChainInfoRefreshRate  time.Duration
	UpgradeRefreshRate    time.Duration
	BlockTimeRefreshRate  time.Duration
	QueryValidators       bool
}

func (c Config) GetProviderOrConsumerHost() string {
	if c.ProviderRPCHost != "" {
		return c.ProviderRPCHost
	}

	return c.RPCHost
}

func (c Config) IsProvider() bool {
	return c.ProviderRPCHost != ""
}
