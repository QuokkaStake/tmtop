package config

import "time"

type Config struct {
	RPCHost     string
	RefreshRate time.Duration
	LogLevel    string
}
