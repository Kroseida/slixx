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
	if p.PacketId() == (&supervisorPacket.StatusUpdate{}).PacketId() {
		return h.HandleStatusUpdate(client, p.(*supervisorPacket.StatusUpdate))
	}
	if p.PacketId() == (&supervisorPacket.DeleteInfo{}).PacketId() {
		return h.HandleDeleteInfo(client, p.(*supervisorPacket.DeleteInfo))
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

func (h *Handler) HandleStatusUpdate(_ protocol.WrappedClient, update *supervisorPacket.StatusUpdate) error {
	return executionService.ApplyExecutionToIndex(
		update.Id,
		update.Kind,
		update.JobId,
		update.Percentage,
		update.StatusType,
		update.Message,
	)
}

func (h *Handler) HandleDeleteInfo(_ protocol.WrappedClient, info *supervisorPacket.DeleteInfo) error {
	return backupService.DeleteBackupFromIndex(info.Id)
}
