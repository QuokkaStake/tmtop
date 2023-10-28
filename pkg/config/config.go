package config

import "time"

type Config struct {
	RPCHost               string
	RefreshRate           time.Duration
	ValidatorsRefreshRate time.Duration
	LogLevel              string
	QueryValidators       bool
}
