package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type Logger interface {
	Info(args ...interface{})
	Debug(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
}

func CreateLogger(loggerMode string, logFile string) *zap.SugaredLogger {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.Encoding = "console"
	if loggerMode == "debug" {
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else if loggerMode == "error" {
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	} else {
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	loggerConfig.OutputPaths = []string{"stdout", logFile}

	logger, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}

	return logger.Sugar()
}
