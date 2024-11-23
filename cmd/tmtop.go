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

func Execute(inputConfig configPkg.InputConfig, args []string) {
	if len(args) == 0 || args[0] == "" {
		inputConfig.RPCHost = "http://localhost:26657"
	} else {
		inputConfig.RPCHost = args[0]
	}

	config, err := configPkg.ParseAndValidateConfig(inputConfig)
	if err != nil {
		panic(err)
	}

	app := pkg.NewApp(config, version)
	app.Start()
	// select {}
}

func main() {
	var config configPkg.InputConfig

	rootCmd := &cobra.Command{
		Use:     "tmtop [RPC host URL]",
		Long:    "Observe the pre-voting status of any Tendermint-based blockchain.",
		Version: version,
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			Execute(config, args)
		},
	}

	rootCmd.PersistentFlags().StringVar(&config.ProviderRPCHost, "provider-rpc-host", "", "Provider chain RPC host URL")
	rootCmd.PersistentFlags().StringVar(&config.ConsumerID, "consumer-id", "", "Consumer ID (not chain ID!)")
	rootCmd.PersistentFlags().DurationVar(&config.RefreshRate, "refresh-rate", time.Second, "Refresh rate")
	rootCmd.PersistentFlags().BoolVar(&config.Verbose, "verbose", false, "Display more debug logs")
	rootCmd.PersistentFlags().BoolVar(&config.DisableEmojis, "disable-emojis", false, "Disable emojis in output")
	rootCmd.PersistentFlags().StringVar(&config.ChainType, "chain-type", "cosmos-rpc", "Chain type. Allowed values are: 'cosmos-rpc', 'cosmos-lcd', 'tendermint'")
	rootCmd.PersistentFlags().DurationVar(&config.ValidatorsRefreshRate, "validators-refresh-rate", time.Minute, "Validators refresh rate")
	rootCmd.PersistentFlags().DurationVar(&config.ChainInfoRefreshRate, "chain-info-refresh-rate", 5*time.Minute, "Chain info refresh rate")
	rootCmd.PersistentFlags().DurationVar(&config.UpgradeRefreshRate, "upgrade-refresh-rate", 30*time.Minute, "Upgrades refresh rate")
	rootCmd.PersistentFlags().DurationVar(&config.BlockTimeRefreshRate, "block-time-refresh-rate", 30*time.Second, "Block time refresh rate")
	rootCmd.PersistentFlags().StringVar(&config.LCDHost, "lcd-host", "", "LCD API host URL")
	rootCmd.PersistentFlags().StringVar(&config.DebugFile, "debug-file", "", "Path to file to write debug info/logs to")
	rootCmd.PersistentFlags().Int64Var(&config.HaltHeight, "halt-height", 0, "Custom halt-height")
	rootCmd.PersistentFlags().Uint64Var(&config.BlocksBehind, "blocks-behind", 1000, "How many blocks behind to check to calculate block time")
	rootCmd.PersistentFlags().StringVar(&config.Timezone, "timezone", "", "Timezone to display dates in")
	rootCmd.PersistentFlags().BoolVar(&config.WithTopologyAPI, "with-topology-api", false, "Enable topology API")
	rootCmd.PersistentFlags().StringVar(&config.TopologyListenAddr, "topology-listen-addr", "0.0.0.0:8080", "The address on which to serve topology API")
	rootCmd.PersistentFlags().StringSliceVar(&config.TopologyHighlightNodes, "topology-highlight-nodes", nil, "The nodes to highlight in the produced topology graph, either by node id, IP address or moniker")

	if err := rootCmd.Execute(); err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Could not start application")
	}
}
