package packet

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	"kroseida.org/slixx/pkg/utils/bytebuf"
)

type ExecuteRestore struct {
	Id       *uuid.UUID `json:"id"`
	JobId    uuid.UUID  `json:"jobId"`
	BackupId uuid.UUID  `json:"backupId"`
}

func (packet *ExecuteRestore) PacketId() int64 {
	return 11
}

func (packet *ExecuteRestore) Protocol() []string {
	return []string{protocol.Supervisor}
}

func (packet *ExecuteRestore) Serialize(buffer *bytebuf.ByteBuffer) error {
	buffer.WriteString(packet.Id.String())
	buffer.WriteString(packet.JobId.String())
	buffer.WriteString(packet.BackupId.String())
	return nil
}

func (packet *ExecuteRestore) Deserialize(buffer *bytebuf.ByteBuffer) error {
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

	backupId, err := uuid.Parse(buffer.ReadString())
	if err != nil {
		return err
	}
	packet.BackupId = backupId

	return nil
}
