package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

func GetDefaultLogger() *zerolog.Logger {
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	return &log
}

type Writer struct {
	io.Writer
	LogChannel chan string
}

func (w Writer) Write(msg []byte) (int, error) {
	w.LogChannel <- string(msg)
	return len(msg), nil
}

func GetLogger(logChannel chan string) *zerolog.Logger {
	writer := zerolog.ConsoleWriter{Out: Writer{LogChannel: logChannel}, NoColor: true}
	log := zerolog.New(writer).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	return &log
}
