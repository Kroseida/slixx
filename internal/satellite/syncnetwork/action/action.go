package action

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/internal/satellite/application"
	"kroseida.org/slixx/internal/satellite/syncnetwork/manager"
	"kroseida.org/slixx/pkg/strategy"
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	"kroseida.org/slixx/pkg/syncnetwork/protocol/supervisor/packet"
	"time"
)

func SyncLogsToSupervisor() {
	// iterate over all connection in server
	for _, connection := range manager.Server.ActiveConnection {
		if connection.Protocol != protocol.Supervisor {
			continue
		}
		syncedLogger := manager.Server.Logger.(*application.SyncedLogger)
		logs := syncedLogger.GetLinesAndClear()
		if len(logs) == 0 {
			continue
		}
		for _, log := range logs {
			log.SatelliteId = *manager.Server.Id
		}

		connection.Send(&packet.SyncLogs{
			Logs: logs,
		})
	}
}

func SendBackupStatusUpdate(id *uuid.UUID, status strategy.BackupStatusUpdate) {
	// iterate over all connection in server
	for _, connection := range manager.Server.ActiveConnection {
		if connection.Protocol != protocol.Supervisor {
			continue
		}

		connection.Send(&packet.BackupStatusUpdate{
			Id:         *id,
			JobId:      *status.JobId,
			Percentage: status.Percentage,
			StatusType: status.StatusType,
			Message:    status.Message,
		})
	}
}

func SendRawBackupInfo(id *uuid.UUID, jobId *uuid.UUID, executionId *uuid.UUID, date time.Time,
	originKind string, destinationKind string, strategy string) {
	// iterate over all connection in server
	for _, connection := range manager.Server.ActiveConnection {
		if connection.Protocol != protocol.Supervisor {
			continue
		}

		connection.Send(&packet.RawBackupInfo{
			Id:              id,
			JobId:           jobId,
			ExecutionId:     executionId,
			CreatedAt:       date,
			OriginKind:      originKind,
			DestinationKind: destinationKind,
			Strategy:        strategy,
		})
	}
}
