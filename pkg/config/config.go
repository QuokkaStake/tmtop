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

	LCDHost string
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
	if !slices.Contains([]string{"cosmos-rpc", "cosmos-lcd", "tendermint"}, c.ChainType) {
		return fmt.Errorf("expected chain-type to be one of 'cosmos-rpc', 'cosmos-lcd', 'tendermint', but got '%s'", c.ChainType)
	}

	if c.IsConsumer() && c.ConsumerChainID == "" {
		return errors.New("chain is consumer, but consumer-chain-id is not set")
	}

	if c.ChainType == "cosmos-lcd" && c.LCDHost == "" {
		return errors.New("chain-type is 'cosmos-lcd', but lcd-node is not set")
	}

	return nil
}
