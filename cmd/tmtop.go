package main

import (
	"github.com/spf13/cobra"
	"main/pkg"
	configPkg "main/pkg/config"
	"main/pkg/logger"
	"time"
)

var (
	version = "unknown"
)

func Execute(config configPkg.Config) {
	app := pkg.NewApp(config, version)
	app.Start()
}

func main() {
	var config configPkg.Config

	rootCmd := &cobra.Command{
		Use:     "tmtop",
		Long:    "Observe the pre-voting status of any Tendermint-based blockchain.",
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			Execute(config)
		},
	}

	rootCmd.PersistentFlags().StringVar(&config.RPCHost, "rpc-host", "http://localhost:26657", "RPC host URL")
	rootCmd.PersistentFlags().StringVar(&config.LogLevel, "log-level", "info", "Log level")
	rootCmd.PersistentFlags().DurationVar(&config.RefreshRate, "refresh-rate", time.Second, "Refresh rate")

	requiredFlags := []string{
		//"rpc-host",
	}

	for _, requiredFlag := range requiredFlags {
		if err := rootCmd.MarkPersistentFlagRequired(requiredFlag); err != nil {
			logger.GetDefaultLogger().Fatal().Str("flag", requiredFlag).Err(err).Msg("Could not set flag as required")
		}
	}

	if err := rootCmd.Execute(); err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Could not start application")
	}
}