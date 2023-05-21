package satellite

import (
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
)

type Handler struct {
}

func (h *Handler) Handle(client protocol.WrappedClient, packet protocol.Packet) error {
	return nil
}
