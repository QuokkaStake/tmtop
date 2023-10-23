package pkg

import (
	"fmt"
	"github.com/rs/zerolog"
	configPkg "main/pkg/config"
	loggerPkg "main/pkg/logger"
)

type App struct {
	Logger  zerolog.Logger
	Version string
	Config  configPkg.Config
}

func NewApp(config configPkg.Config, version string) *App {
	logger := loggerPkg.GetLogger(config.LogLevel).
		With().
		Str("component", "app_manager").
		Logger()

	return &App{
		Logger:  logger,
		Version: version,
		Config:  config,
	}
}

func (a *App) Start() {
	fmt.Printf("config: %+v\n", a.Config)
}
