package main

import (
	"kroseida.org/slixx/internal/common"
	"kroseida.org/slixx/internal/supervisor/application"
	"kroseida.org/slixx/internal/supervisor/datasource"
	"kroseida.org/slixx/internal/supervisor/graphql"
	satelliteService "kroseida.org/slixx/internal/supervisor/service/satellite"
	"kroseida.org/slixx/pkg/utils"
	"os"
)

var SETTINGS = "supervisor.settings.json"

func main() {
	err := utils.LoadSettings(SETTINGS, &application.CurrentSettings, &application.DefaultSettings)
	if err != nil {
		panic(err)
	}
	application.Logger = utils.CreateLogger(application.CurrentSettings.Logger.Mode, "supervisor.log")
	application.Logger.Info("Starting Slixx supervisor v" + common.CurrentVersion)

	application.Logger.Info("Initializing database connection")
	err = datasource.Connect()
	if err != nil {
		application.Logger.Error("Failed to initialize database connection: ", err)
		os.Exit(1)
		return
	}

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
