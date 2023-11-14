package action

import (
	"github.com/google/uuid"
	"github.com/samsarahq/thunder/graphql"
	"kroseida.org/slixx/internal/supervisor/application"
	"kroseida.org/slixx/internal/supervisor/datasource"
	syncnetworkClients "kroseida.org/slixx/internal/supervisor/syncnetwork/clients"
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
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
		if client.Client.Protocol != protocol.Supervisor {
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
		if client.Client.Protocol != protocol.Supervisor {
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
		if client.Client.Protocol != protocol.Supervisor {
			continue
		}
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

func SendRequestBackupSync(id *uuid.UUID) error {
	var hasSent bool
	for clientId, client := range syncnetworkClients.List {
		if id != nil && clientId != *id {
			continue
		}
		if client.Client.Protocol != protocol.Supervisor {
			continue
		}
		err := client.Client.Send(&supervisorPacket.RequestBackupSync{})
		if err != nil {
			return err
		}
		hasSent = true
	}
	if !hasSent {
		return graphql.NewSafeError("satellite not connected")
	}
	return nil
}