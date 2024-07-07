package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func Init(logLevelStr string) {
	config := zap.NewProductionConfig()

	var level zapcore.Level
	if err := level.Set(logLevelStr); err != nil {
		level = zapcore.InfoLevel
	}

	config.Level = zap.NewAtomicLevelAt(level)

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	Logger = logger
}
