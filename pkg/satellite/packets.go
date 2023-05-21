package satellite

import (
	"kroseida.org/slixx/pkg/satellite/protocol"
	"kroseida.org/slixx/pkg/satellite/protocol/handshake/packet"
)

var PACKETS = map[int64]protocol.Packet{
	(&packet.Handshake{}).PacketId():          &packet.Handshake{},
	(&packet.ConnectionDenied{}).PacketId():   &packet.ConnectionDenied{},
	(&packet.ConnectionAccepted{}).PacketId(): &packet.ConnectionAccepted{},
}
