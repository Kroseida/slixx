package supervisor

import (
	"kroseida.org/slixx/internal/supervisor/service/satellitelogservice"
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	supervisorPacket "kroseida.org/slixx/pkg/syncnetwork/protocol/supervisor/packet"
)

type Handler struct {
}

func (h *Handler) Handle(client protocol.WrappedClient, p protocol.Packet) error {
	if p.PacketId() == (&supervisorPacket.SyncLogs{}).PacketId() {
		return h.HandleSyncLogs(client, p.(*supervisorPacket.SyncLogs))
	}
	return nil
}

func (h *Handler) HandleSyncLogs(_ protocol.WrappedClient, logs *supervisorPacket.SyncLogs) error {
	return satellitelogservice.Create(logs.Logs)
}
