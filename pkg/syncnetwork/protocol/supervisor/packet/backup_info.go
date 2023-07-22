package packet

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	"kroseida.org/slixx/pkg/utils/bytebuf"
	"time"
)

type RawBackupInfo struct {
	Id    *uuid.UUID `json:"id"`
	JobId *uuid.UUID `json:"job_id"`
	Date  time.Time  `json:"date"`
}

func (packet *RawBackupInfo) PacketId() int64 {
	return 9
}

func (packet *RawBackupInfo) Protocol() []string {
	return []string{protocol.Supervisor}
}

func (packet *RawBackupInfo) Serialize(buffer *bytebuf.ByteBuffer) error {
	buffer.WriteString(packet.Id.String())
	buffer.WriteString(packet.JobId.String())
	buffer.WriteString(packet.Date.Format(time.RFC3339))
	return nil
}

func (packet *RawBackupInfo) Deserialize(buffer *bytebuf.ByteBuffer) error {
	id, err := uuid.Parse(buffer.ReadString())
	if err != nil {
		return err
	}
	packet.Id = &id

	jobId, err := uuid.Parse(buffer.ReadString())
	if err != nil {
		return err
	}
	packet.JobId = &jobId

	date, err := time.Parse(time.RFC3339, buffer.ReadString())
	if err != nil {
		return err
	}
	packet.Date = date

	return nil
}
