package packet

import (
	"encoding/json"
	"kroseida.org/slixx/pkg/model"
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	"kroseida.org/slixx/pkg/utils/bytebuf"
)

type SyncStorage struct {
	Storages []*model.Storage `json:"storages"`
}

func (packet *SyncStorage) PacketId() int64 {
	return 3
}

func (packet *SyncStorage) Protocol() []string {
	return []string{protocol.Supervisor}
}

func (packet *SyncStorage) Serialize(buffer *bytebuf.ByteBuffer) error {
	storagesJson, err := json.Marshal(packet.Storages)
	if err != nil {
		return err
	}

	buffer.Write(storagesJson)

	return nil
}

func (packet *SyncStorage) Deserialize(buffer *bytebuf.ByteBuffer) error {
	var storages []*model.Storage

	err := json.Unmarshal(buffer.Read(), &storages)
	if err != nil {
		return err
	}

	packet.Storages = storages

	return nil
}
