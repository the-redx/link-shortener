package utils

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func getLevelLogger(level string) zapcore.Level {
	if level == "DEBUG" {
		return zap.DebugLevel
	}

	return zap.ErrorLevel
}

func init() {
	appLogLevel := os.Getenv("APP_LOG_LEVEL")
	appEnv := os.Getenv("APP_ENV")

	zapConfig := zap.NewProductionConfig()
	zapConfig.EncoderConfig = zap.NewProductionEncoderConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(getLevelLogger(appLogLevel))
	zapConfig.Development = appEnv == "development"
	zapConfig.Encoding = "json"
	zapConfig.OutputPaths = []string{"stdout"}
	zapConfig.ErrorOutputPaths = []string{"stderr"}

	logger, err := zapConfig.Build()
	if err != nil {
		panic(err)
	}

	Logger = logger.Sugar()
}
