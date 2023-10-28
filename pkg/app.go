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
	}
}

func (a *App) Start() {
	go a.GoRefreshConsensus()

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
	consensus, validators, chainValidators, err := a.Aggregator.GetData()
	if err != nil {
		a.Logger.Error().Err(err).Msg("Error getting data")
		return
	}

	state, err := types.StateFromTendermintResponse(consensus, validators, chainValidators)
	if err != nil {
		a.Logger.Error().Err(err).Msg("Error converting data")
		return
	}
	// fmt.Printf("parsed %+v\n", state)

	a.DisplayWrapper.SetState(state)
}
