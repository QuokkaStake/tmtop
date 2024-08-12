package logger

import (
	"io"
	configPkg "main/pkg/config"
	"os"

	"github.com/rs/zerolog"
)

func GetDefaultLogger() *zerolog.Logger {
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	return &log
}

type Writer struct {
	io.Writer
	DebugFile  *os.File
	LogChannel chan string
}

func NewWriter(logChannel chan string, config *configPkg.Config) Writer {
	writer := Writer{
		LogChannel: logChannel,
	}

	if config.DebugFile != "" {
		debugFile, err := os.OpenFile(config.DebugFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			panic(err)
		}

		writer.DebugFile = debugFile
	}

	return writer
}

func (w Writer) Write(msg []byte) (int, error) {
	w.LogChannel <- string(msg)

	if w.DebugFile != nil {
		if _, err := w.DebugFile.Write(msg); err != nil {
			return 0, err
		}
	}

	return len(msg), nil
}

func GetLogger(logChannel chan string, config *configPkg.Config) *zerolog.Logger {
	writer := zerolog.ConsoleWriter{
		Out:     NewWriter(logChannel, config),
		NoColor: true,
	}
	log := zerolog.New(writer).With().Timestamp().Logger()

	if config.Verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	}

	return &log
}
