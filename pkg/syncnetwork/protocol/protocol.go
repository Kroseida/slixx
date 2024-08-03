package protocol

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"kroseida.org/slixx/pkg/utils/bytebuf"
	"reflect"
)

var Supervisor = "SUPERVISOR"
var Satellite = "SATELLITE"
var Handshake = "HANDSHAKE"

type Packet interface {
	PacketId() int64
	Protocol() []string
	Serialize(buffer *bytebuf.ByteBuffer) error
	Deserialize(buffer *bytebuf.ByteBuffer) error
}

func SendPacket(writer bufio.Writer, packet Packet) error {
	buffer := bytebuf.Get()
	defer bytebuf.Put(buffer)

	// Send Packet id as second
	buffer.WriteInt64(packet.PacketId())

	err := packet.Serialize(buffer)
	if err != nil {
		return err
	}

	length := bytebuf.Get()
	defer bytebuf.Put(length)

	// Send Packet Length as first
	length.WriteInt64(int64(buffer.Len()))

	data := append(length.Bytes(), buffer.Bytes()...)

	_, err = writer.Write(data)
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
	// Read Packet Length as first
	length, err := binary.ReadVarint(reader)
	if err != nil {
		return nil, err
	}

	buffer := bytebuf.Get()
	defer bytebuf.Put(buffer)

	data := make([]byte, length)

	_, err = reader.Read(data)
	if err != nil {
		return nil, err
	}

	buffer.SetBytes(data)

	if err != nil {
		return nil, err
	}
	packetId := buffer.ReadInt64()

	packetType := packetRegistry[packetId]
	if packetType == nil {
		return nil, fmt.Errorf("unknown packet *id: %d", packetId)
	}
	packet := reflect.New(reflect.TypeOf(packetType).Elem()).Interface().(Packet)

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
	IsConnected() bool
}
