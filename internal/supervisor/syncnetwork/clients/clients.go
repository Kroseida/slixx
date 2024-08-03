package clients

import (
	"github.com/google/uuid"
	supervisorProtocol "kroseida.org/slixx/internal/supervisor/syncnetwork/protocol"
)

var List = make(map[uuid.UUID]*supervisorProtocol.WrappedClient)

func ListConnected() map[uuid.UUID]*supervisorProtocol.WrappedClient {
	connected := make(map[uuid.UUID]*supervisorProtocol.WrappedClient)
	for id, client := range List {
		if client.Client.Connected {
			connected[id] = client
		}
	}
	return connected
}

func GetConnectedClient(id uuid.UUID) *supervisorProtocol.WrappedClient {
	client := List[id]
	if client == nil || !client.Client.Connected {
		return nil
	}
	return client
}
