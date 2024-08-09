package config

import (
	"errors"
	"fmt"
	"time"
)

type ChainType string

const (
	ChainTypeCosmosRPC  ChainType = "cosmos-rpc"
	ChainTypeCosmosLCD  ChainType = "cosmos-lcd"
	ChainTypeTendermint ChainType = "tendermint"
)

func (t *ChainType) String() string {
	return string(*t)
}

func (t *ChainType) Set(v string) error {
	switch v {
	case "cosmos-rpc", "":
		*t = ChainTypeCosmosRPC
		return nil
	case "cosmos-lcd":
		*t = ChainTypeCosmosLCD
		return nil
	case "tendermint":
		*t = ChainTypeTendermint
		return nil
	}

	return fmt.Errorf(
		"expected chain-type to be one of 'cosmos-rpc', 'cosmos-lcd', 'tendermint', but got '%s'",
		v,
	)
}
func (t *ChainType) Type() string {
	return "ChainType"
}

type Config struct {
	RPCHost               string
	ProviderRPCHost       string
	ConsumerChainID       string
	RefreshRate           time.Duration
	ValidatorsRefreshRate time.Duration
	ChainInfoRefreshRate  time.Duration
	UpgradeRefreshRate    time.Duration
	BlockTimeRefreshRate  time.Duration
	ChainType             ChainType
	Verbose               bool
	DisableEmojis         bool
	DebugFile             string
	HaltHeight            int64

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
	if c.IsConsumer() && c.ConsumerChainID == "" {
		return errors.New("chain is consumer, but consumer-chain-id is not set")
	}

	if c.ChainType == ChainTypeCosmosLCD && c.LCDHost == "" {
		return errors.New("chain-type is 'cosmos-lcd', but lcd-host is not set")
	}

	return nil
}
