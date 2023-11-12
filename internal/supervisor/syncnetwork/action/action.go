package action

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/internal/supervisor/application"
	"kroseida.org/slixx/internal/supervisor/datasource"
	syncnetworkClients "kroseida.org/slixx/internal/supervisor/syncnetwork/clients"
	supervisorPacket "kroseida.org/slixx/pkg/syncnetwork/protocol/supervisor/packet"
)

func SyncStorages(id *uuid.UUID) {
	storages, err := datasource.StorageProvider.List()

	if err != nil {
		application.Logger.Error("Failed to load storages from database for sync: ", err)
		return
	}

	for clientId, client := range syncnetworkClients.List {
		if id != nil && clientId != *id {
			continue
		}
		client.Client.Send(&supervisorPacket.SyncStorage{
			Storages: storages,
		})
	}
}

func SyncJobs(id *uuid.UUID) {
	jobs, err := datasource.JobProvider.List()

	if err != nil {
		application.Logger.Error("Failed to load jobs from database for sync: ", err)
		return
	}

	for clientId, client := range syncnetworkClients.List {
		if id != nil && clientId != *id {
			continue
		}
		err := client.Client.Send(&supervisorPacket.SyncJob{
			Jobs: jobs,
		})
		if err != nil {
			application.Logger.Error("Failed to send jobs to client: ", err)
		}
	}
}

func SendExecuteBackup(jobId uuid.UUID) (*uuid.UUID, error) {
	id := uuid.New()
	for _, client := range syncnetworkClients.List {
		err := client.Client.Send(&supervisorPacket.ExecuteBackup{
			Id:    &id,
			JobId: jobId,
		})
		if err != nil {
			return nil, err
		}
	}
	return &id, nil
}

func SendRequestBackupSync(satelliteId uuid.UUID) error {
	for _, client := range syncnetworkClients.List {
		err := client.Client.Send(&supervisorPacket.RequestBackupSync{
			SatelliteId: satelliteId,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
