package syncnetwork

import (
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	handshakePacket "kroseida.org/slixx/pkg/syncnetwork/protocol/handshake/packet"
	supervisorPacket "kroseida.org/slixx/pkg/syncnetwork/protocol/supervisor/packet"
)

var PACKETS = map[int64]protocol.Packet{
	(&handshakePacket.Handshake{}).PacketId():          &handshakePacket.Handshake{},
	(&handshakePacket.ConnectionDenied{}).PacketId():   &handshakePacket.ConnectionDenied{},
	(&handshakePacket.ConnectionAccepted{}).PacketId(): &handshakePacket.ConnectionAccepted{},
	(&supervisorPacket.SyncStorage{}).PacketId():       &supervisorPacket.SyncStorage{},
	(&supervisorPacket.SyncJob{}).PacketId():           &supervisorPacket.SyncJob{},
	(&supervisorPacket.SyncLogs{}).PacketId():          &supervisorPacket.SyncLogs{},
	(&supervisorPacket.ApplySupervisor{}).PacketId():   &supervisorPacket.ApplySupervisor{},
	(&supervisorPacket.ExecuteBackup{}).PacketId():     &supervisorPacket.ExecuteBackup{},
	(&supervisorPacket.StatusUpdate{}).PacketId():      &supervisorPacket.StatusUpdate{},
	(&supervisorPacket.RawBackupInfo{}).PacketId():     &supervisorPacket.RawBackupInfo{},
	(&supervisorPacket.RequestResync{}).PacketId():     &supervisorPacket.RequestResync{},
	(&supervisorPacket.ExecuteRestore{}).PacketId():    &supervisorPacket.ExecuteRestore{},
	(&supervisorPacket.SyncJobSchedule{}).PacketId():   &supervisorPacket.SyncJobSchedule{},
	(&supervisorPacket.DeleteBackup{}).PacketId():      &supervisorPacket.DeleteBackup{},
	(&supervisorPacket.DeleteInfo{}).PacketId():        &supervisorPacket.DeleteInfo{},
}
