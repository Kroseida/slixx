package protocol

import (
	"bufio"
	"fmt"
	"kroseida.org/slixx/pkg/utils/bytebuf"
)

var SupervisorProtocol = "SUPERVISOR"
var SatelliteProtocol = "SATELLITE"
var HandshakeProtocol = "HANDSHAKE"

type Packet interface {
	PacketId() int64
	Protocol() []string
	Serialize(buffer *bytebuf.ByteBuffer) error
	Deserialize(buffer *bytebuf.ByteBuffer) error
}

func SendPacket(writer bufio.Writer, packet Packet) error {
	buffer := bytebuf.Get()
	defer bytebuf.Put(buffer)
	buffer.WriteInt64(packet.PacketId())

	err := packet.Serialize(buffer)
	if err != nil {
		return err
	}

	_, err = writer.Write(append(buffer.Bytes(), 0x04))
	if err != nil {
		return err
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func ReadPacket(reader *bufio.Reader, packetRegistry map[int64]Packet) (Packet, error) {
	buffer := bytebuf.Get()
	defer bytebuf.Put(buffer)

	data, err := reader.ReadBytes(0x04)
	buffer.SetBytes(data)

	if err != nil {
		return nil, err
	}
	packetId := buffer.ReadInt64()

	packet := packetRegistry[packetId]
	if packet == nil {
		return nil, fmt.Errorf("unknown packet id: %d", packetId)
	}
	err = packet.Deserialize(buffer)
	if err != nil {
		return nil, err
	}
	return packet, nil
}

type Handler interface {
	Handle(client WrappedClient, packet Packet) error
}

type WrappedClient interface {
	Send(packet Packet) error
}
