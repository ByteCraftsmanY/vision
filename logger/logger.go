package logger

import (
	"go.uber.org/zap"
	"vision/config"
)

func Init(config *config.Configurations) *zap.Logger {
	log := zap.Must(zap.NewProduction())
	zap.ReplaceGlobals(log)
	return log
}
