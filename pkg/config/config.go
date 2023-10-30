package config

import "time"

type Config struct {
	RPCHost               string
	ProviderRPCHost       string
	RefreshRate           time.Duration
	ValidatorsRefreshRate time.Duration
	ChainInfoRefreshRate  time.Duration
	QueryValidators       bool
}

func (c Config) GetProviderOrConsumerHost() string {
	if c.ProviderRPCHost != "" {
		return c.ProviderRPCHost
	}

	return c.RPCHost
}
