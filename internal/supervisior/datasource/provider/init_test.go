package provider_test

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"kroseida.org/slixx/internal/supervisior/application"
	"kroseida.org/slixx/internal/supervisior/datasource"
	"os"
	"time"
)

func setupSuite() func() {
	databaseName := "slixx.db"
	os.Remove(databaseName)

	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.Encoding = "console"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	logger, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}

	application.Logger = logger.Sugar()
	datasource.Connect()

	return func() {
		datasource.Close()
		os.Remove(databaseName)
	}
}
