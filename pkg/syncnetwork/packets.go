package syncnetwork

import (
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	"kroseida.org/slixx/pkg/syncnetwork/protocol/handshake/packet"
)

var PACKETS = map[int64]protocol.Packet{
	(&packet.Handshake{}).PacketId():          &packet.Handshake{},
	(&packet.ConnectionDenied{}).PacketId():   &packet.ConnectionDenied{},
	(&packet.ConnectionAccepted{}).PacketId(): &packet.ConnectionAccepted{},
}
