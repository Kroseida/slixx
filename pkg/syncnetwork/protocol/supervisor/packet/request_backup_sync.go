package packet

import (
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	"kroseida.org/slixx/pkg/utils/bytebuf"
)

type RequestBackupSync struct {
}

func (packet *RequestBackupSync) PacketId() int64 {
	return 10
}

func (packet *RequestBackupSync) Protocol() []string {
	return []string{protocol.Supervisor}
}

func (packet *RequestBackupSync) Serialize(buffer *bytebuf.ByteBuffer) error {
	return nil
}

func (packet *RequestBackupSync) Deserialize(buffer *bytebuf.ByteBuffer) error {
	return nil
}
