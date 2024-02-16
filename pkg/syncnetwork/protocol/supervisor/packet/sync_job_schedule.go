package packet

import (
	"encoding/json"
	"kroseida.org/slixx/pkg/model"
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	"kroseida.org/slixx/pkg/utils/bytebuf"
)

type SyncJobSchedule struct {
	Schedules []*model.JobSchedule `json:"jobSchedules"`
}

func (packet *SyncJobSchedule) PacketId() int64 {
	return 12
}

func (packet *SyncJobSchedule) Protocol() []string {
	return []string{protocol.Supervisor}
}

func (packet *SyncJobSchedule) Serialize(buffer *bytebuf.ByteBuffer) error {
	storagesJson, err := json.Marshal(packet.Schedules)
	if err != nil {
		return err
	}

	buffer.Write(storagesJson)

	return nil
}

func (packet *SyncJobSchedule) Deserialize(buffer *bytebuf.ByteBuffer) error {
	var schedules []*model.JobSchedule

	err := json.Unmarshal(buffer.Read(), &schedules)
	if err != nil {
		return err
	}

	packet.Schedules = schedules

	return nil
}
