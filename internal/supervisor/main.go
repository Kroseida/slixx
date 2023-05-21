package main

import (
	"kroseida.org/slixx/internal/supervisor/application"
	"kroseida.org/slixx/internal/supervisor/datasource"
	"kroseida.org/slixx/internal/supervisor/graphql"
	"kroseida.org/slixx/internal/supervisor/syncnetwork"
	"kroseida.org/slixx/pkg/utils"
	"os"
)

var SETTINGS = "supervisor.settings.json"

func main() {
	err := utils.LoadSettings(SETTINGS, &application.CurrentSettings, &application.DefaultSettings)
	if err != nil {
		panic(err)
	}
	application.Logger = utils.CreateLogger(application.CurrentSettings.Logger.Mode)
	application.Logger.Info("Starting Slixx supervisor")

	application.Logger.Info("Initializing database connection")
	err = datasource.Connect()
	if err != nil {
		application.Logger.Errorw("Failed to initialize database connection", "error", err)
		os.Exit(1)
		return
	}

	// We check for new satellites every ... minutes
	go syncnetwork.Watchdog()

	application.Logger.Info("Initializing GraphQL API Server on " + application.CurrentSettings.Http.BindAddress)
	err = graphql.Start()
	if err != nil {
		application.Logger.Errorw("Failed to initialize GraphQL API Server", "error", err)
		os.Exit(1)
		return
	}
}
