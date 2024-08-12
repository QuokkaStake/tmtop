package config

import (
	"errors"
	"fmt"
	"time"
)

type InputConfig struct {
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
	DebugFile             string
	HaltHeight            int64
	LCDHost               string
}

type ChainType string

const (
	ChainTypeCosmosRPC  ChainType = "cosmos-rpc"
	ChainTypeCosmosLCD  ChainType = "cosmos-lcd"
	ChainTypeTendermint ChainType = "tendermint"
)

func (t *ChainType) String() string {
	return string(*t)
}

func ParseChainType(v string) (ChainType, error) {
	switch v {
	case "cosmos-rpc", "":
		return ChainTypeCosmosRPC, nil
	case "cosmos-lcd":
		return ChainTypeCosmosLCD, nil
	case "tendermint":
		return ChainTypeTendermint, nil
	}

	return "", fmt.Errorf(
		"expected chain-type to be one of 'cosmos-rpc', 'cosmos-lcd', 'tendermint', but got '%s'",
		v,
	)
}

func ParseAndValidateConfig(input InputConfig) (*Config, error) {
	chainType, err := ParseChainType(input.ChainType)
	if err != nil {
		return nil, err
	}

	if input.ProviderRPCHost != "" && input.ConsumerChainID == "" {
		return nil, errors.New("chain is consumer, but consumer-chain-id is not set")
	}

	if chainType == ChainTypeCosmosLCD && input.LCDHost == "" {
		return nil, errors.New("chain-type is 'cosmos-lcd', but lcd-host is not set")
	}

	config := &Config{
		RPCHost:               input.RPCHost,
		ProviderRPCHost:       input.ProviderRPCHost,
		ConsumerChainID:       input.ConsumerChainID,
		RefreshRate:           input.RefreshRate,
		ValidatorsRefreshRate: input.ValidatorsRefreshRate,
		ChainInfoRefreshRate:  input.ChainInfoRefreshRate,
		UpgradeRefreshRate:    input.UpgradeRefreshRate,
		BlockTimeRefreshRate:  input.BlockTimeRefreshRate,
		ChainType:             chainType,
		Verbose:               input.Verbose,
		DisableEmojis:         input.DisableEmojis,
		DebugFile:             input.DebugFile,
		HaltHeight:            input.HaltHeight,
		LCDHost:               input.LCDHost,
	}

	return config, nil
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
	LCDHost               string
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
