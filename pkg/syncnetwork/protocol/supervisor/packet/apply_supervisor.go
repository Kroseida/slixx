package packet

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	"kroseida.org/slixx/pkg/utils/bytebuf"
)

type ApplySupervisor struct {
	Id uuid.UUID `json:"uuid"`
}

func (packet *ApplySupervisor) PacketId() int64 {
	return 6
}

func (packet *ApplySupervisor) Protocol() []string {
	return []string{protocol.Supervisor}
}

func (packet *ApplySupervisor) Serialize(buffer *bytebuf.ByteBuffer) error {
	buffer.WriteString(packet.Id.String())
	return nil
}

func (packet *ApplySupervisor) Deserialize(buffer *bytebuf.ByteBuffer) error {
	id, err := uuid.Parse(buffer.ReadString())
	if err != nil {
		return err
	}
	packet.Id = id
	return nil
}
