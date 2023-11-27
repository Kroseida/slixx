package packet

import (
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	"kroseida.org/slixx/pkg/utils/bytebuf"
)

type RequestResync struct {
}

func (packet *RequestResync) PacketId() int64 {
	return 10
}

func (packet *RequestResync) Protocol() []string {
	return []string{protocol.Supervisor}
}

func (packet *RequestResync) Serialize(buffer *bytebuf.ByteBuffer) error {
	return nil
}

func (packet *RequestResync) Deserialize(buffer *bytebuf.ByteBuffer) error {
	return nil
}
