package supervisor

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/internal/satellite/backup"
	"kroseida.org/slixx/internal/satellite/syncdata"
	"kroseida.org/slixx/internal/satellite/syncnetwork/action"
	"kroseida.org/slixx/pkg/model"
	strategyRegistry "kroseida.org/slixx/pkg/strategy"
	syncnetworkBase "kroseida.org/slixx/pkg/syncnetwork"
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	supervisorPacket "kroseida.org/slixx/pkg/syncnetwork/protocol/supervisor/packet"
	"kroseida.org/slixx/pkg/utils"
	"strconv"
)

type Handler struct {
}

func (h *Handler) Handle(client protocol.WrappedClient, packet protocol.Packet) error {
	if packet.PacketId() == (&supervisorPacket.SyncStorage{}).PacketId() {
		return h.HandleSyncStorage(client, packet.(*supervisorPacket.SyncStorage))
	}
	if packet.PacketId() == (&supervisorPacket.SyncJob{}).PacketId() {
		return h.HandleSyncJob(client, packet.(*supervisorPacket.SyncJob))
	}
	if packet.PacketId() == (&supervisorPacket.ApplySupervisor{}).PacketId() {
		return h.HandleApplySupervisor(client, packet.(*supervisorPacket.ApplySupervisor))
	}
	if packet.PacketId() == (&supervisorPacket.ExecuteBackup{}).PacketId() {
		return h.HandleExecuteBackup(client, packet.(*supervisorPacket.ExecuteBackup))
	}
	return nil
}

func (h *Handler) HandleSyncStorage(client protocol.WrappedClient, storage *supervisorPacket.SyncStorage) error {
	c := client.(*syncnetworkBase.ConnectedClient)

	syncdata.Container.Storages = map[uuid.UUID]*model.Storage{}

	for _, storage := range storage.Storages {
		syncdata.Container.Storages[storage.Id] = storage
	}
	syncdata.GenerateCache()

	c.Server.Logger.Info("Synced " + strconv.Itoa(len(syncdata.Container.Storages)) + " storages from supervisor")
	return nil
}

func (h *Handler) HandleSyncJob(client protocol.WrappedClient, job *supervisorPacket.SyncJob) error {
	c := client.(*syncnetworkBase.ConnectedClient)

	jobIdsBeforeSync := make([]uuid.UUID, len(syncdata.Container.Jobs))
	for _, job := range syncdata.Container.Jobs {
		jobIdsBeforeSync = append(jobIdsBeforeSync, job.Id)
	}

	syncdata.Container.Jobs = map[uuid.UUID]*model.Job{}

	for _, job := range job.Jobs {
		syncdata.Container.Jobs[job.Id] = job
	}
	syncdata.GenerateCache()

	c.Server.Logger.Info("Synced " + strconv.Itoa(len(syncdata.Container.Jobs)) + " jobs from supervisor")

	jobIdsAfterSync := make([]uuid.UUID, len(syncdata.Container.Jobs))
	for _, job := range syncdata.Container.Jobs {
		jobIdsAfterSync = append(jobIdsAfterSync, job.Id)
	}

	newJobIds := make([]uuid.UUID, 0)
	for _, jobId := range jobIdsAfterSync {
		if !utils.ContainsUUID(jobIdsBeforeSync, jobId) {
			newJobIds = append(newJobIds, jobId)
		}
	}
	backup.SendBackupInfos(newJobIds)

	return nil
}

func (h *Handler) HandleApplySupervisor(client protocol.WrappedClient, supervisor *supervisorPacket.ApplySupervisor) error {
	c := client.(*syncnetworkBase.ConnectedClient)
	c.Server.Id = &supervisor.Id
	return nil
}

func (h *Handler) HandleExecuteBackup(_ protocol.WrappedClient, execute *supervisorPacket.ExecuteBackup) error {
	go func() {
		err := backup.Execute(execute.Id, execute.JobId)
		if err != nil {
			action.SendBackupStatusUpdate(execute.Id, strategyRegistry.BackupStatusUpdate{
				JobId:      &execute.JobId,
				Percentage: 0,
				StatusType: "FAILED",
				Message:    err.Error(),
			})
		}
	}()
	return nil
}
