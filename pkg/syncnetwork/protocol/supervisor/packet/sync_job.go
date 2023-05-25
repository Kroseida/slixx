package packet

import (
	"encoding/json"
	"kroseida.org/slixx/pkg/model"
	"kroseida.org/slixx/pkg/utils/bytebuf"
)

type SyncJob struct {
	Jobs []*model.Job
}

func (packet *SyncJob) PacketId() int64 {
	return 4
}

func (packet *SyncJob) Protocol() []string {
	return []string{"supervisor"}
}

func (packet *SyncJob) Serialize(buffer *bytebuf.ByteBuffer) error {
	jobsJson, err := json.Marshal(packet.Jobs)
	if err != nil {
		return err
	}

	buffer.Write(jobsJson)

	return nil
}

func (packet *SyncJob) Deserialize(buffer *bytebuf.ByteBuffer) error {
	var jobs []*model.Job

	err := json.Unmarshal(buffer.Read(), &jobs)
	if err != nil {
		return err
	}

	packet.Jobs = jobs

	return nil
}
