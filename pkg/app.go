package pkg

import (
	"fmt"
	"main/pkg/aggregator"
	configPkg "main/pkg/config"
	"main/pkg/display"
	tmhttp "main/pkg/http"
	loggerPkg "main/pkg/logger"
	"main/pkg/topology"
	"main/pkg/types"
	"sync"
	"time"

	"github.com/brynbellomy/go-utils"
	"github.com/rs/zerolog"
)

type App struct {
	Logger         zerolog.Logger
	Version        string
	Config         *configPkg.Config
	Aggregator     *aggregator.Aggregator
	DisplayWrapper *display.Wrapper
	State          *types.State
	LogChannel     chan string

	mbRPCURLs        *utils.Mailbox[types.RPC]
	rpcURLsLastFetch map[string]time.Time

	PauseChannel chan bool
	IsPaused     bool
}

func NewApp(config *configPkg.Config, version string) *App {
	logChannel := make(chan string, 1000)
	pauseChannel := make(chan bool)

	state := types.NewState(config.RPCHost)

	logger := loggerPkg.GetLogger(logChannel, config).
		With().
		Str("component", "app_manager").
		Logger()

	return &App{
		Logger:           logger,
		Version:          version,
		Config:           config,
		Aggregator:       aggregator.NewAggregator(config, state, logger),
		DisplayWrapper:   display.NewWrapper(config, state, logger, pauseChannel, version),
		State:            state,
		LogChannel:       logChannel,
		mbRPCURLs:        utils.NewMailbox[types.RPC](1000),
		rpcURLsLastFetch: make(map[string]time.Time),
		PauseChannel:     pauseChannel,
		IsPaused:         false,
	}
}

func (a *App) Start() {
	if a.Config.WithTopologyAPI {
		go a.ServeTopology()
		topology.LogChannel = a.LogChannel
	}

	go a.CrawlRPCURLs()

	go a.GoRefreshConsensus()
	go a.GoRefreshValidators()
	go a.GoRefreshChainInfo()
	go a.GoRefreshUpgrade()
	go a.GoRefreshBlockTime()
	go a.GoRefreshNetInfo()
	go a.DisplayLogs()
	go a.ListenForPause()

	a.DisplayWrapper.Start()
}

func (a *App) ServeTopology() {
	_ = tmhttp.NewServer(
		a.Config.TopologyListenAddr,
		topology.WithHTTPTopologyAPI(a.State),
		topology.WithHTTPPeersAPI(a.State),
		topology.WithFrontendStaticAssets(),
	).Serve()
}

func (a *App) CrawlRPCURLs() {
	a.fetchNewPeers(a.Config.RPCHost)
	timer := time.NewTimer(15 * time.Second)

	for {
		select {
		case <-a.mbRPCURLs.Notify():
			var wg sync.WaitGroup
			for _, rpc := range a.mbRPCURLs.RetrieveAll() {
				if lastFetch, ok := a.rpcURLsLastFetch[rpc.URL]; ok && time.Now().Sub(lastFetch) < 15*time.Second {
					continue
				}
				a.rpcURLsLastFetch[rpc.URL] = time.Now()
				a.State.AddKnownRPC(rpc)

				rpc := rpc
				wg.Add(1)
				go func() {
					defer wg.Done()

					a.fetchNewPeers(rpc.URL)
				}()
			}
			wg.Wait()

		case <-timer.C:
			for _, rpc := range a.State.KnownRPCs().Iter() {
				if time.Since(a.rpcURLsLastFetch[rpc.URL]) >= 15*time.Second {
					a.mbRPCURLs.Deliver(rpc)
				}
			}
		}

	}
}

func (a *App) fetchNewPeers(rpcURL string) {
	netInfo, err := a.Aggregator.GetNetInfo(rpcURL)
	if err != nil {
		a.LogChannel <- fmt.Sprintf("error getting net_info from %s: %v", rpcURL, err)
		return
	}

	a.State.AddRPCPeers(rpcURL, netInfo.Peers)
	for _, peer := range netInfo.Peers {
		a.mbRPCURLs.Deliver(types.NewRPCFromPeer(peer))
	}
}

func (a *App) GoRefreshConsensus() {
	defer a.HandlePanic()

	a.RefreshConsensus()

	ticker := time.NewTicker(a.Config.RefreshRate)
	done := make(chan bool)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			a.RefreshConsensus()
		}
	}
}

func (a *App) RefreshConsensus() {
	if a.IsPaused {
		return
	}

	consensus, validators, err := a.Aggregator.GetData()
	a.State.SetConsensusStateError(err)
	if err != nil {
		a.Logger.Error().Err(err).Msg("Error getting consensus data")
		a.DisplayWrapper.SetState(a.State)
		return
	}

	err = a.State.SetTendermintResponse(consensus, validators)
	a.State.SetConsensusStateError(err)
	if err != nil {
		a.Logger.Error().Err(err).Msg("Error converting data")
		a.DisplayWrapper.SetState(a.State)
		return
	}

	a.State.SetConsensusStateError(err)
	a.DisplayWrapper.SetState(a.State)
}

func (a *App) GoRefreshValidators() {
	defer a.HandlePanic()

	a.RefreshValidators()

	ticker := time.NewTicker(a.Config.ValidatorsRefreshRate)
	done := make(chan bool)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			a.RefreshValidators()
		}
	}
}

func (a *App) RefreshValidators() {
	if a.IsPaused {
		return
	}

	chainValidators, err := a.Aggregator.GetChainValidators()
	if err != nil {
		a.DisplayWrapper.SetState(a.State)
		a.Logger.Error().Err(err).Msg("Error getting chain validators")
		return
	}

	a.State.SetChainValidators(chainValidators)
	a.DisplayWrapper.SetState(a.State)
}

func (a *App) GoRefreshChainInfo() {
	defer a.HandlePanic()

	a.RefreshChainInfo()

	ticker := time.NewTicker(a.Config.ChainInfoRefreshRate)
	done := make(chan bool)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			a.RefreshChainInfo()
		}
	}
}

func (a *App) RefreshChainInfo() {
	if a.IsPaused {
		return
	}

	chainInfo, err := a.Aggregator.GetChainInfo()
	if err != nil {
		a.Logger.Error().Err(err).Msg("Error getting chain validators")
		a.State.SetChainInfoError(err)
		a.DisplayWrapper.SetState(a.State)
		return
	}

	a.State.SetChainInfo(&chainInfo.Result)
	a.State.SetChainInfoError(err)
	a.DisplayWrapper.SetState(a.State)
}

func (a *App) GoRefreshUpgrade() {
	defer a.HandlePanic()

	a.RefreshUpgrade()

	ticker := time.NewTicker(a.Config.UpgradeRefreshRate)
	done := make(chan bool)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			a.RefreshUpgrade()
		}
	}
}

func (a *App) RefreshUpgrade() {
	if a.IsPaused {
		return
	}

	if a.Config.HaltHeight > 0 {
		upgrade := &types.Upgrade{
			Name:   "halt-height upgrade",
			Height: a.Config.HaltHeight,
		}

		a.State.SetUpgrade(upgrade)
		a.DisplayWrapper.SetState(a.State)
		return
	}

	upgrade, err := a.Aggregator.GetUpgrade()
	if err != nil {
		a.Logger.Error().Err(err).Msg("Error getting upgrade")
		a.State.SetUpgradePlanError(err)
		a.DisplayWrapper.SetState(a.State)
		return
	}

	a.State.SetUpgrade(upgrade)
	a.State.SetUpgradePlanError(err)
	a.DisplayWrapper.SetState(a.State)
}

func (a *App) GoRefreshBlockTime() {
	defer a.HandlePanic()

	a.RefreshBlockTime()

	ticker := time.NewTicker(a.Config.BlockTimeRefreshRate)
	done := make(chan bool)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			a.RefreshBlockTime()
		}
	}
}

func (a *App) RefreshBlockTime() {
	if a.IsPaused {
		return
	}

	blockTime, err := a.Aggregator.GetBlockTime()
	if err != nil {
		a.Logger.Error().Err(err).Msg("Error getting block time")
		return
	}

	a.State.SetBlockTime(blockTime)
	a.DisplayWrapper.SetState(a.State)
}

func (a *App) GoRefreshNetInfo() {
	defer a.HandlePanic()

	a.RefreshNetInfo()

	ticker := time.NewTicker(a.Config.RefreshRate)
	done := make(chan bool)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			a.RefreshNetInfo()
		}
	}
}

func (a *App) RefreshNetInfo() {
	if a.IsPaused {
		return
	}

	netInfo, err := a.Aggregator.GetNetInfo(a.State.CurrentRPC().URL)
	if err != nil {
		a.Logger.Error().Err(err).Msg("Error getting netInfo")
		return
	}

	a.State.SetNetInfo(netInfo)
	a.DisplayWrapper.SetState(a.State)
}

func (a *App) DisplayLogs() {
	for {
		logString := <-a.LogChannel
		a.DisplayWrapper.DebugText(logString)
	}
}

func (a *App) ListenForPause() {
	for {
		paused := <-a.PauseChannel
		a.IsPaused = paused
	}
}

func (a *App) HandlePanic() {
	if r := recover(); r != nil {
		a.DisplayWrapper.App.Stop()
		panic(r)
	}
}
