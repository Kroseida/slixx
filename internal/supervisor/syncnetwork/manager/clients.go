package manager

import (
	"github.com/google/uuid"
	supervisorProtocol "kroseida.org/slixx/internal/supervisor/syncnetwork/protocol"
)

var Clients = make(map[uuid.UUID]*supervisorProtocol.WrappedClient)
