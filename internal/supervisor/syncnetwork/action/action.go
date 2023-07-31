package action

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/internal/supervisor/application"
	"kroseida.org/slixx/internal/supervisor/datasource"
	syncnetworkClients "kroseida.org/slixx/internal/supervisor/syncnetwork/clients"
	supervisorPacket "kroseida.org/slixx/pkg/syncnetwork/protocol/supervisor/packet"
)

func SyncStorages() {
	storages, err := datasource.StorageProvider.List()

	if err != nil {
		application.Logger.Error("Failed to load storages from database for sync: ", err)
		return
	}

	for _, client := range syncnetworkClients.List {
		client.Client.Send(&supervisorPacket.SyncStorage{
			Storages: storages,
		})
	}
}

func SyncJobs() {
	jobs, err := datasource.JobProvider.List()

	if err != nil {
		application.Logger.Error("Failed to load jobs from database for sync: ", err)
		return
	}

	for _, client := range syncnetworkClients.List {
		err := client.Client.Send(&supervisorPacket.SyncJob{
			Jobs: jobs,
		})
		if err != nil {
			application.Logger.Error("Failed to send jobs to client: ", err)
		}
	}
}

func SendExecuteBackup(jobId uuid.UUID) uuid.UUID {
	id := uuid.New()
	for _, client := range syncnetworkClients.List {
		client.Client.Send(&supervisorPacket.ExecuteBackup{
			Id:    &id,
			JobId: jobId,
		})
	}
	return id
}
