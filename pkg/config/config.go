package config

import "time"

type Config struct {
	RPCHost               string
	RefreshRate           time.Duration
	ValidatorsRefreshRate time.Duration
	ChainInfoRefreshRate  time.Duration
	QueryValidators       bool
}
