package packet

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	"kroseida.org/slixx/pkg/utils/bytebuf"
	"time"
)

type RawBackupInfo struct {
	Id              *uuid.UUID `json:"id"`
	JobId           *uuid.UUID `json:"job_id"`
	ExecutionId     uuid.UUID  `json:"execution_id"`
	CreatedAt       time.Time  `json:"created_at"`
	OriginKind      string     `json:"origin_kind"`
	DestinationKind string     `json:"destination_kind"`
	Strategy        string     `json:"strategy"`
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
	buffer.WriteString(packet.ExecutionId.String())
	buffer.WriteString(packet.CreatedAt.Format(time.RFC3339))
	buffer.WriteString(packet.OriginKind)
	buffer.WriteString(packet.DestinationKind)
	buffer.WriteString(packet.Strategy)
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

	executionId, err := uuid.Parse(buffer.ReadString())
	if err != nil {
		return err
	}
	packet.ExecutionId = executionId

	date, err := time.Parse(time.RFC3339, buffer.ReadString())
	if err != nil {
		return err
	}
	packet.CreatedAt = date

	packet.OriginKind = buffer.ReadString()
	packet.DestinationKind = buffer.ReadString()
	packet.Strategy = buffer.ReadString()

	return nil
}
