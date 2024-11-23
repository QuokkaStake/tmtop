package config

import (
	"errors"
	"fmt"
	"time"
)

type InputConfig struct {
	RPCHost               string
	ProviderRPCHost       string
	ConsumerID            string
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
	BlocksBehind          uint64
	LCDHost               string
	Timezone              string
	WithTopologyAPI       bool
	TopologyListenAddr    string
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

	if input.ProviderRPCHost != "" && input.ConsumerID == "" {
		return nil, errors.New("chain is consumer, but consumer-id is not set")
	}

	if chainType == ChainTypeCosmosLCD && input.LCDHost == "" {
		return nil, errors.New("chain-type is 'cosmos-lcd', but lcd-host is not set")
	}

	if input.BlocksBehind <= 0 {
		return nil, errors.New("cannot run with a negative blocks-behind")
	}

	timezone := time.Local //nolint:gosmopolitan // local timezone is expected here

	if input.Timezone != "" {
		parsedTimezone, err := time.LoadLocation(input.Timezone)
		if err != nil {
			return nil, err
		}

		timezone = parsedTimezone
	}

	config := &Config{
		RPCHost:               input.RPCHost,
		ProviderRPCHost:       input.ProviderRPCHost,
		ConsumerID:            input.ConsumerID,
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
		BlocksBehind:          input.BlocksBehind,
		LCDHost:               input.LCDHost,
		Timezone:              timezone,
		WithTopologyAPI:       input.WithTopologyAPI,
		TopologyListenAddr:    input.TopologyListenAddr,
	}

	return config, nil
}

type Config struct {
	RPCHost               string
	ProviderRPCHost       string
	ConsumerID            string
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
	BlocksBehind          uint64
	LCDHost               string
	Timezone              *time.Location
	WithTopologyAPI       bool
	TopologyListenAddr    string
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
