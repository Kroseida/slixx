package main

import (
	"kroseida.org/slixx/internal/satellite/application"
	"kroseida.org/slixx/internal/satellite/syncnetwork"
	"kroseida.org/slixx/pkg/utils"
)

var SETTINGS = "satellite.settings.json"

func main() {
	err := utils.LoadSettings(SETTINGS, &application.CurrentSettings, &application.DefaultSettings)
	if err != nil {
		panic(err)
	}
	application.Logger = utils.CreateLogger(application.CurrentSettings.Logger.Mode)
	application.Logger.Info("Starting Slixx satellite")

	application.Logger.Info("Listening sync network on " + application.CurrentSettings.Satellite.Network.BindAddress)
	syncnetwork.Listen()
}
