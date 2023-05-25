package syncnetwork

import (
	"kroseida.org/slixx/internal/supervisor/application"
	"kroseida.org/slixx/internal/supervisor/datasource"
	supervisorPacket "kroseida.org/slixx/pkg/syncnetwork/protocol/supervisor/packet"
)

func SyncStorages() {
	storages, err := datasource.StorageProvider.GetStorages()

	if err != nil {
		application.Logger.Error("Failed to load storages from database for sync: ", err)
		return
	}

	for _, client := range clients {
		client.Client.Send(&supervisorPacket.SyncStorage{
			Storages: storages,
		})
	}
}

func SyncJobs() {
	jobs, err := datasource.JobProvider.GetJobs()

	if err != nil {
		application.Logger.Error("Failed to load jobs from database for sync: ", err)
		return
	}

	for _, client := range clients {
		client.Client.Send(&supervisorPacket.SyncJob{
			Jobs: jobs,
		})
	}
}
