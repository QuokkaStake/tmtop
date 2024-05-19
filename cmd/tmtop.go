package main

import (
	"main/pkg"
	configPkg "main/pkg/config"
	"main/pkg/logger"
	"time"

	"github.com/spf13/cobra"
)

var (
	version = "unknown"
)

func Execute(config configPkg.Config, args []string) {
	if len(args) == 0 || args[0] == "" {
		config.RPCHost = "http://localhost:26657"
	} else {
		config.RPCHost = args[0]
	}

	app := pkg.NewApp(config, version)
	app.Start()
}

func main() {
	var config configPkg.Config

	rootCmd := &cobra.Command{
		Use:     "tmtop [RPC host URL]",
		Long:    "Observe the pre-voting status of any Tendermint-based blockchain.",
		Version: version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return config.Validate()
		},
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			Execute(config, args)
		},
	}

	rootCmd.PersistentFlags().StringVar(&config.ProviderRPCHost, "provider-rpc-host", "", "Provider chain RPC host URL")
	rootCmd.PersistentFlags().StringVar(&config.ConsumerChainID, "consumer-chain-id", "", "Consumer chain ID")
	rootCmd.PersistentFlags().DurationVar(&config.RefreshRate, "refresh-rate", time.Second, "Refresh rate")
	rootCmd.PersistentFlags().BoolVar(&config.Verbose, "verbose", false, "Display more debug logs")
	rootCmd.PersistentFlags().BoolVar(&config.DisableEmojis, "disable-emojis", false, "Disable emojis in output")
	rootCmd.PersistentFlags().Var(&config.ChainType, "chain-type", "Chain type. Allowed values are: 'cosmos-rpc', 'cosmos-lcd', 'tendermint'")
	rootCmd.PersistentFlags().DurationVar(&config.ValidatorsRefreshRate, "validators-refresh-rate", time.Minute, "Validators refresh rate")
	rootCmd.PersistentFlags().DurationVar(&config.ChainInfoRefreshRate, "chain-info-refresh-rate", 5*time.Minute, "Chain info refresh rate")
	rootCmd.PersistentFlags().DurationVar(&config.UpgradeRefreshRate, "upgrade-refresh-rate", 30*time.Minute, "Upgrades refresh rate")
	rootCmd.PersistentFlags().DurationVar(&config.BlockTimeRefreshRate, "block-time-refresh-rate", 30*time.Second, "Block time refresh rate")
	rootCmd.PersistentFlags().StringVar(&config.LCDHost, "lcd-host", "", "LCD API host URL")
	rootCmd.PersistentFlags().StringVar(&config.DebugFile, "debug-file", "", "Path to file to write debug info/logs to")

	if err := rootCmd.Execute(); err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Could not start application")
	}
}
