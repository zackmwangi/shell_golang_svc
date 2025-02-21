package applogger

import (
	"go.uber.org/zap"
)

type (
	AppLogger struct {
		Logger *zap.Logger
	}
)

func InitAppLogger() *zap.Logger {

	logger, err := zap.NewProduction()

	if err != nil {
		panic("could not init app logger")
	}

	return logger
}
