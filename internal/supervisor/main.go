package main

import (
	"kroseida.org/slixx/internal/common"
	"kroseida.org/slixx/internal/supervisor/application"
	"kroseida.org/slixx/internal/supervisor/datasource"
	"kroseida.org/slixx/internal/supervisor/graphql"
	"kroseida.org/slixx/internal/supervisor/service/execution"
	satelliteService "kroseida.org/slixx/internal/supervisor/service/satellite"
	"kroseida.org/slixx/internal/supervisor/service/satellitelog"
	"kroseida.org/slixx/pkg/utils"
	"kroseida.org/slixx/pkg/utils/fileutils"
	"os"
)

var SETTINGS = "config/supervisor.settings.json"

func main() {
	if !fileutils.FileExists("log") {
		os.Mkdir("log", 0755)
	}
	if !fileutils.FileExists("config") {
		os.Mkdir("config", 0755)
	}

	err := utils.LoadSettings(SETTINGS, &application.CurrentSettings, &application.DefaultSettings)
	if err != nil {
		panic(err)
	}

	application.Logger = utils.CreateLogger(application.CurrentSettings.Logger.Mode, "log/supervisor.log")
	application.Logger.Info("Starting Slixx supervisor v" + common.CurrentVersion)

	application.Logger.Info("Initializing database connection")
	err = datasource.Connect()
	if err != nil {
		application.Logger.Error("Failed to initialize database connection: ", err)
		os.Exit(1)
		return
	}

	// We start the log cleanup job
	go satellitelog.StartCleanupJob()
	go execution.StartTimeoutDetector()

	// We check for new satellites every ... minutes
	go satelliteService.StartWatchdog()

	application.Logger.Info("Initializing GraphQL API Server on " + application.CurrentSettings.Http.BindAddress)
	err = graphql.Start()
	if err != nil {
		application.Logger.Error("Failed to initialize GraphQL API Server", err)
		os.Exit(1)
		return
	}
}
