package clients

import (
	"github.com/google/uuid"
	supervisorProtocol "kroseida.org/slixx/internal/supervisor/syncnetwork/protocol"
)

var List = make(map[uuid.UUID]*supervisorProtocol.WrappedClient)
