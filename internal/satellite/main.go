package main

import (
	"kroseida.org/slixx/internal/common"
	"kroseida.org/slixx/internal/satellite/application"
	"kroseida.org/slixx/internal/satellite/syncdata"
	"kroseida.org/slixx/internal/satellite/syncnetwork"
	"kroseida.org/slixx/pkg/model"
	"kroseida.org/slixx/pkg/utils"
	"kroseida.org/slixx/pkg/utils/fileutils"
	"os"
)

var SETTINGS = "config/satellite.settings.json"

func main() {
	if !fileutils.FileExists("data") {
		os.Mkdir("data", 0755)
	}

	err := utils.LoadSettings(SETTINGS, &application.CurrentSettings, &application.DefaultSettings)
	if err != nil {
		panic(err)
	}
	application.Logger = &application.SyncedLogger{
		Logger:      utils.CreateLogger(application.CurrentSettings.Logger.Mode, "log/satellite.log"),
		CachedLines: []*model.SatelliteLogEntry{},
	}

	application.Logger.Info("Starting Slixx satellite v" + common.CurrentVersion)
	application.Logger.Info("Loading cache from disk")
	err = syncdata.LoadCache()
	if err != nil {
		application.Logger.Error("Failed to load cache from disk", err)
	}

	application.Logger.Info("Listening sync network on " + application.CurrentSettings.Satellite.Network.BindAddress)
	go syncnetwork.SyncLoop()
	syncnetwork.Listen()
}
