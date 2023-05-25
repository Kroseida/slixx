package packet

import (
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	"kroseida.org/slixx/pkg/utils/bytebuf"
)

type Handshake struct {
	Id             string `json:"id"`
	Version        string `json:"version"`
	Token          string `json:"token"`
	TargetProtocol string `json:"targetProtocol"`
}

func (packet *Handshake) PacketId() int64 {
	return 0
}

func (packet *Handshake) Protocol() []string {
	return []string{protocol.Handshake}
}

func (packet *Handshake) Serialize(buffer *bytebuf.ByteBuffer) error {
	buffer.WriteString(packet.Version)
	buffer.WriteString(packet.Token)
	buffer.WriteString(packet.TargetProtocol)
	return nil
}

func (packet *Handshake) Deserialize(buffer *bytebuf.ByteBuffer) error {
	packet.Version = buffer.ReadString()
	packet.Token = buffer.ReadString()
	packet.TargetProtocol = buffer.ReadString()
	return nil
}
