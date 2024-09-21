package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
)

type Module struct {
	Logger zerolog.Logger
}

func NewModule() *Module {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	if viper.Get("ENV") != "prod" {
		// global pretty logging 效果
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		// 如果只想要個別創建的實例有 pretty logging 效果
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: zerolog.TimeFormatUnix}
		logger = zerolog.New(consoleWriter).With().Timestamp().Logger()
	}

	return &Module{
		Logger: logger,
	}
}
