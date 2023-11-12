package packet

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	"kroseida.org/slixx/pkg/utils/bytebuf"
)

type RequestBackupSync struct {
	SatelliteId uuid.UUID `json:"satelliteId"`
}

func (packet *RequestBackupSync) PacketId() int64 {
	return 10
}

func (packet *RequestBackupSync) Protocol() []string {
	return []string{protocol.Supervisor}
}

func (packet *RequestBackupSync) Serialize(buffer *bytebuf.ByteBuffer) error {
	buffer.WriteString(packet.SatelliteId.String())
	return nil
}

func (packet *RequestBackupSync) Deserialize(buffer *bytebuf.ByteBuffer) error {
	satelliteId, err := uuid.Parse(buffer.ReadString())
	if err != nil {
		return err
	}
	packet.SatelliteId = satelliteId
	return nil
}
