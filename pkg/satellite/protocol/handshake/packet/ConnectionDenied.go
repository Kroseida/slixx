package packet

import (
	"kroseida.org/slixx/pkg/satellite/protocol"
	"kroseida.org/slixx/pkg/utils/bytebuf"
)

type ConnectionDenied struct{}

func (packet *ConnectionDenied) PacketId() int64 {
	return 1
}

func (packet *ConnectionDenied) Protocol() []string {
	return []string{protocol.HandshakeProtocol, protocol.SupervisorProtocol, protocol.SatelliteProtocol}
}

func (packet *ConnectionDenied) Serialize(buffer *bytebuf.ByteBuffer) error {
	return nil
}

func (packet *ConnectionDenied) Deserialize(buffer *bytebuf.ByteBuffer) error {
	return nil
}
