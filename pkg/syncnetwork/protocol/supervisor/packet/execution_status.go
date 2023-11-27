package packet

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	"kroseida.org/slixx/pkg/utils/bytebuf"
)

type ExecutionStatusUpdate struct {
	Id         uuid.UUID `json:"id"`
	JobId      uuid.UUID `json:"jobId"`
	Percentage float64   `json:"percentage"`
	StatusType string    `json:"statusType"`
	Message    string    `json:"message"`
}

func (packet *ExecutionStatusUpdate) PacketId() int64 {
	return 7
}

func (packet *ExecutionStatusUpdate) Protocol() []string {
	return []string{protocol.Supervisor}
}

func (packet *ExecutionStatusUpdate) Serialize(buffer *bytebuf.ByteBuffer) error {
	buffer.WriteString(packet.Id.String())
	buffer.WriteString(packet.JobId.String())
	buffer.WriteFloat64(packet.Percentage)
	buffer.WriteString(packet.StatusType)
	buffer.WriteString(packet.Message)
	return nil
}

func (packet *ExecutionStatusUpdate) Deserialize(buffer *bytebuf.ByteBuffer) error {
	id, err := uuid.Parse(buffer.ReadString())
	if err != nil {
		return err
	}
	jobId, err := uuid.Parse(buffer.ReadString())
	if err != nil {
		return err
	}
	packet.Id = id
	packet.JobId = jobId
	packet.Percentage = buffer.ReadFloat64()
	packet.StatusType = buffer.ReadString()
	packet.Message = buffer.ReadString()
	return nil
}
