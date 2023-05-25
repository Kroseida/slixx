package syncnetwork

import (
	"kroseida.org/slixx/internal/satellite/application"
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	"kroseida.org/slixx/pkg/syncnetwork/protocol/supervisor/packet"
)

func SyncLogsToSupervisor() {
	if server == nil {
		return
	}
	if server.Id == nil {
		return
	}

	// iterate over all connection in server
	for _, connection := range server.ActiveConnection {
		if connection.Protocol != protocol.Supervisor {
			continue
		}
		syncedLogger := server.Logger.(*application.SyncedLogger)
		logs := syncedLogger.GetLinesAndClear()
		if len(logs) == 0 {
			continue
		}
		for _, log := range logs {
			log.SatelliteId = *server.Id
		}

		connection.Send(&packet.SyncLogs{
			Logs: logs,
		})
	}
}
