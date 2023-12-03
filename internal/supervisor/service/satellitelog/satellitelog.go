package satellitelog

import (
	"kroseida.org/slixx/internal/supervisor/application"
	"kroseida.org/slixx/internal/supervisor/datasource"
	"kroseida.org/slixx/pkg/model"
	"time"
)

func StartCleanupJob() {
	if !application.CurrentSettings.LogSync.Active {
		return
	}
	for {
		application.Logger.Info("Deleting synced logs older than ", application.CurrentSettings.LogSync.LogRetention, " hours")
		err := datasource.SatelliteProvider.
			DeleteLogsOlderThan(time.Now().Add(-time.Hour * time.Duration(application.CurrentSettings.LogSync.LogRetention)))

		if err != nil {
			application.Logger.Error("Failed to delete old logs: ", err)
		}

		time.Sleep(time.Hour * time.Duration(application.CurrentSettings.LogSync.CheckInterval))
	}
}

func Create(logs []*model.SatelliteLogEntry) error {
	if !application.CurrentSettings.LogSync.Active {
		return nil
	}
	return datasource.SatelliteProvider.ApplyLogs(logs)
}
