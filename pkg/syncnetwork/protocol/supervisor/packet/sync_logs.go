package packet

import (
	"encoding/json"
	"kroseida.org/slixx/pkg/model"
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	"kroseida.org/slixx/pkg/utils/bytebuf"
)

type SyncLogs struct {
	Logs []*model.SatelliteLogEntry `json:"logs"`
}

func (packet *SyncLogs) PacketId() int64 {
	return 5
}

func (packet *SyncLogs) Protocol() []string {
	return []string{protocol.Supervisor}
}

func (packet *SyncLogs) Serialize(buffer *bytebuf.ByteBuffer) error {
	logsJson, err := json.Marshal(packet.Logs)
	if err != nil {
		return err
	}

	buffer.Write(logsJson)

	return nil
}

func (packet *SyncLogs) Deserialize(buffer *bytebuf.ByteBuffer) error {
	var logs []*model.SatelliteLogEntry

	err := json.Unmarshal(buffer.Read(), &logs)
	if err != nil {
		return err
	}

	packet.Logs = logs

	return nil
}
