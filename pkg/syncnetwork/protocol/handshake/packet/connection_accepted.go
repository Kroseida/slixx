package packet

import (
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	"kroseida.org/slixx/pkg/utils/bytebuf"
)

type ConnectionAccepted struct{}

func (packet *ConnectionAccepted) PacketId() int64 {
	return 2
}

func (packet *ConnectionAccepted) Protocol() []string {
	return []string{protocol.Handshake, protocol.Supervisor, protocol.Satellite}
}

func (packet *ConnectionAccepted) Serialize(buffer *bytebuf.ByteBuffer) error {
	return nil
}

func (packet *ConnectionAccepted) Deserialize(buffer *bytebuf.ByteBuffer) error {
	return nil
}
