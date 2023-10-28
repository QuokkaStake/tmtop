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
}

func NewApp(config configPkg.Config, version string) *App {
	logger := loggerPkg.GetLogger(config.LogLevel).
		With().
		Str("component", "app_manager").
		Logger()

	return &App{
		Logger:         logger,
		Version:        version,
		Config:         config,
		Aggregator:     aggregator.NewAggregator(config, logger),
		DisplayWrapper: display.NewWrapper(logger),
		State:          types.NewState(),
	}
}

func (a *App) Start() {
	go a.GoRefreshConsensus()
	go a.GoRefreshValidators()

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
	consensus, validators, err := a.Aggregator.GetData()
	if err != nil {
		a.Logger.Error().Err(err).Msg("Error getting consensus data")
		return
	}

	err = a.State.SetTendermintResponse(consensus, validators)
	if err != nil {
		a.Logger.Error().Err(err).Msg("Error converting data")
		return
	}
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
	chainValidators, err := a.Aggregator.GetChainValidators()
	if err != nil {
		a.Logger.Error().Err(err).Msg("Error getting chain validators")
		return
	}

	a.State.SetChainValidators(chainValidators)
	a.DisplayWrapper.SetState(a.State)
}
