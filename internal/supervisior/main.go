package main

import (
	"encoding/json"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"kroseida.org/slixx/internal/supervisior/application"
	"kroseida.org/slixx/internal/supervisior/datasource"
	"kroseida.org/slixx/internal/supervisior/graphql"
	"kroseida.org/slixx/pkg/utils"
	"os"
	"time"
)

var SETTINGS = "settings.json"

func main() {
	err := LoadSettings()
	if err != nil {
		panic(err)
	}
	CreateLogger()
	application.Logger.Info("Starting Slixx supervisior")

	application.Logger.Info("Initializing database connection")
	err = datasource.Connect()
	if err != nil {
		application.Logger.Errorw("Failed to initialize database connection", "error", err)
		os.Exit(1)
		return
	}
	
	application.Logger.Info("Initializing GraphQL API Server")
	err = graphql.Start()
	if err != nil {
		application.Logger.Errorw("Failed to initialize GraphQL API Server", "error", err)
		os.Exit(1)
		return
	}
}

func CreateLogger() {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.Encoding = "console"
	if application.CurrentSettings.Logger.Mode == "debug" {
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else if application.CurrentSettings.Logger.Mode == "error" {
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	} else {
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	logger, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}

	application.Logger = logger.Sugar()
}

func LoadSettings() error {
	var err error
	if !utils.FileExists(SETTINGS) {
		err = CreateSettings()
	}
	if err != nil {
		return err
	}

	content, err := os.ReadFile(SETTINGS)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &application.CurrentSettings)
	if err != nil {
		return err
	}

	return nil
}

func CreateSettings() error {
	content, err := json.MarshalIndent(&application.DefaultSettings, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(SETTINGS, content, 0644)
	if err != nil {
		return err
	}

	return nil
}
