package packet

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	"kroseida.org/slixx/pkg/utils/bytebuf"
)

type DeleteBackup struct {
	Id       uuid.UUID `json:"id"`
	JobId    uuid.UUID `json:"jobId"`
	BackupId uuid.UUID `json:"backupId"`
}

func (packet *DeleteBackup) PacketId() int64 {
	return 13
}

func (packet *DeleteBackup) Protocol() []string {
	return []string{protocol.Supervisor}
}

func (packet *DeleteBackup) Serialize(buffer *bytebuf.ByteBuffer) error {
	buffer.WriteString(packet.Id.String())
	buffer.WriteString(packet.JobId.String())
	buffer.WriteString(packet.BackupId.String())
	return nil
}

func (packet *DeleteBackup) Deserialize(buffer *bytebuf.ByteBuffer) error {
	id, err := uuid.Parse(buffer.ReadString())
	if err != nil {
		return err
	}
	packet.Id = id

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
