package pkg

import (
	"github.com/rs/zerolog"
	configPkg "main/pkg/config"
	loggerPkg "main/pkg/logger"
	"main/pkg/tendermint"
	"main/pkg/types"
	"main/pkg/view_wrapper"
	"time"
)

type App struct {
	Logger         zerolog.Logger
	Version        string
	Config         configPkg.Config
	RPC            *tendermint.RPC
	DisplayWrapper *view_wrapper.Wrapper
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
		RPC:            tendermint.NewRPC(config, logger),
		DisplayWrapper: view_wrapper.NewWrapper(logger),
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
	consensus, validators, err := a.RPC.GetConsensusStateAndValidators()
	if err != nil {
		a.Logger.Fatal().Err(err).Msg("Error getting consensus state")
		return
	}

	renderInfo := types.RenderIntoFromTendermintResponse(consensus, validators)
	//fmt.Printf("parsed %+v\n", renderInfo)

	a.DisplayWrapper.SetInfo(renderInfo)
}
