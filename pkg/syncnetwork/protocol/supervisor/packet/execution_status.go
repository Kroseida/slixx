package packet

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	"kroseida.org/slixx/pkg/utils/bytebuf"
)

type StatusUpdate struct {
	Id         uuid.UUID `json:"id"`
	Kind       string    `json:"kind"`
	JobId      uuid.UUID `json:"jobId"`
	Percentage float64   `json:"percentage"`
	StatusType string    `json:"statusType"`
	Message    string    `json:"message"`
}

func (packet *StatusUpdate) PacketId() int64 {
	return 7
}

func (packet *StatusUpdate) Protocol() []string {
	return []string{protocol.Supervisor}
}

func (packet *StatusUpdate) Serialize(buffer *bytebuf.ByteBuffer) error {
	buffer.WriteString(packet.Id.String())
	buffer.WriteString(packet.Kind)
	buffer.WriteString(packet.JobId.String())
	buffer.WriteFloat64(packet.Percentage)
	buffer.WriteString(packet.StatusType)
	buffer.WriteString(packet.Message)
	return nil
}

func (packet *StatusUpdate) Deserialize(buffer *bytebuf.ByteBuffer) error {
	id, err := uuid.Parse(buffer.ReadString())
	if err != nil {
		return err
	}

	packet.Kind = buffer.ReadString()

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
