package pkg

import (
	"main/pkg/aggregator"
	configPkg "main/pkg/config"
	"main/pkg/display"
	loggerPkg "main/pkg/logger"
	"main/pkg/types"
	"time"

	"github.com/rs/zerolog"
)

type App struct {
	Logger         zerolog.Logger
	Version        string
	Config         configPkg.Config
	Aggregator     *aggregator.Aggregator
	DisplayWrapper *display.Wrapper
	State          *types.State
	LogChannel     chan string

	PauseChannel chan bool
	IsPaused     bool
}

func NewApp(config configPkg.Config, version string) *App {
	logChannel := make(chan string)
	pauseChannel := make(chan bool)

	logger := loggerPkg.GetLogger(logChannel, config).
		With().
		Str("component", "app_manager").
		Logger()

	return &App{
		Logger:         logger,
		Version:        version,
		Config:         config,
		Aggregator:     aggregator.NewAggregator(config, logger),
		DisplayWrapper: display.NewWrapper(config, logger, pauseChannel, version),
		State:          types.NewState(),
		LogChannel:     logChannel,
		PauseChannel:   pauseChannel,
		IsPaused:       false,
	}
}

func (a *App) Start() {
	go a.GoRefreshConsensus()
	go a.GoRefreshValidators()
	go a.GoRefreshChainInfo()
	go a.GoRefreshUpgrade()
	go a.GoRefreshBlockTime()
	go a.DisplayLogs()
	go a.ListenForPause()

	a.DisplayWrapper.Start()
}

func (a *App) GoRefreshConsensus() {
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
