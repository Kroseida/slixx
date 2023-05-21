package supervisor

import (
	"kroseida.org/slixx/pkg/satellite/protocol"
)

type Handler struct {
}

func (h *Handler) Handle(client protocol.WrappedClient, packet protocol.Packet) error {
	return nil
}
