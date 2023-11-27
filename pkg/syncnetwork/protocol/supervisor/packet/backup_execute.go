package packet

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	"kroseida.org/slixx/pkg/utils/bytebuf"
)

type ExecuteBackup struct {
	Id    *uuid.UUID `json:"id"`
	JobId uuid.UUID  `json:"jobId"`
}

func (packet *ExecuteBackup) PacketId() int64 {
	return 8
}

func (packet *ExecuteBackup) Protocol() []string {
	return []string{protocol.Supervisor}
}

func (packet *ExecuteBackup) Serialize(buffer *bytebuf.ByteBuffer) error {
	buffer.WriteString(packet.Id.String())
	buffer.WriteString(packet.JobId.String())
	return nil
}

func (packet *ExecuteBackup) Deserialize(buffer *bytebuf.ByteBuffer) error {
	id, err := uuid.Parse(buffer.ReadString())
	if err != nil {
		return err
	}
	packet.Id = &id

	jobId, err := uuid.Parse(buffer.ReadString())
	if err != nil {
		return err
	}
	packet.JobId = jobId

	return nil
}
