package config

import (
	"errors"
	"fmt"
	"slices"
	"time"
)

type Config struct {
	RPCHost               string
	ProviderRPCHost       string
	ConsumerChainID       string
	RefreshRate           time.Duration
	ValidatorsRefreshRate time.Duration
	ChainInfoRefreshRate  time.Duration
	UpgradeRefreshRate    time.Duration
	BlockTimeRefreshRate  time.Duration
	ChainType             string
	Verbose               bool
	DisableEmojis         bool
}

func (c Config) GetProviderOrConsumerHost() string {
	if c.ProviderRPCHost != "" {
		return c.ProviderRPCHost
	}

	return c.RPCHost
}

func (c Config) IsConsumer() bool {
	return c.ProviderRPCHost != ""
}

func (c Config) Validate() error {
	if !slices.Contains([]string{"cosmos", "tendermint"}, c.ChainType) {
		return fmt.Errorf("expected chain-type to be one of 'cosmos', 'tendermint', but got '%s'", c.ChainType)
	}

	if c.IsConsumer() && c.ConsumerChainID == "" {
		return errors.New("chain is consumer, but consumer-chain-id is not set")
	}

	return nil
}
