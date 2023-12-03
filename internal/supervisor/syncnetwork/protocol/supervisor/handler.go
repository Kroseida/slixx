package supervisor

import (
	backupService "kroseida.org/slixx/internal/supervisor/service/backup"
	executionService "kroseida.org/slixx/internal/supervisor/service/execution"
	satellitelogService "kroseida.org/slixx/internal/supervisor/service/satellitelog"
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	supervisorPacket "kroseida.org/slixx/pkg/syncnetwork/protocol/supervisor/packet"
)

type Handler struct {
}

func (h *Handler) Handle(client protocol.WrappedClient, p protocol.Packet) error {
	if p.PacketId() == (&supervisorPacket.SyncLogs{}).PacketId() {
		return h.HandleSyncLogs(client, p.(*supervisorPacket.SyncLogs))
	}
	if p.PacketId() == (&supervisorPacket.RawBackupInfo{}).PacketId() {
		return h.HandleRawBackupInfo(client, p.(*supervisorPacket.RawBackupInfo))
	}
	if p.PacketId() == (&supervisorPacket.ExecutionStatusUpdate{}).PacketId() {
		return h.HandleExecutionStatusUpdate(client, p.(*supervisorPacket.ExecutionStatusUpdate))
	}
	return nil
}

func (h *Handler) HandleSyncLogs(_ protocol.WrappedClient, logs *supervisorPacket.SyncLogs) error {
	return satellitelogService.Create(logs.Logs)
}

func (h *Handler) HandleRawBackupInfo(_ protocol.WrappedClient, info *supervisorPacket.RawBackupInfo) error {
	return backupService.ApplyBackupToIndex(
		*info.Id,
		*info.JobId,
		info.ExecutionId,
		info.CreatedAt,
		info.OriginKind,
		info.DestinationKind,
		info.Strategy,
	)
}

func (h *Handler) HandleExecutionStatusUpdate(_ protocol.WrappedClient, update *supervisorPacket.ExecutionStatusUpdate) error {
	return executionService.ApplyExecutionToIndex(
		update.Id,
		update.Kind,
		update.JobId,
		update.Percentage,
		update.StatusType,
		update.Message,
	)
}
