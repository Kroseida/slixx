package handshake

import (
	"kroseida.org/slixx/pkg/satellite"
	"kroseida.org/slixx/pkg/satellite/protocol"
	"kroseida.org/slixx/pkg/satellite/protocol/handshake/packet"
	"os"
)

type ClientHandler struct {
}

func (handler *ClientHandler) Handle(client protocol.WrappedClient, p protocol.Packet) error {
	if p.PacketId() == (&packet.ConnectionAccepted{}).PacketId() {
		return handler.HandleConnectionAccepted(client, p.(*packet.ConnectionAccepted))
	}
	if p.PacketId() == (&packet.ConnectionDenied{}).PacketId() {
		return handler.HandleConnectionDenied(client, p.(*packet.ConnectionDenied))
	}

	return nil
}

func (handler *ClientHandler) HandleConnectionAccepted(client protocol.WrappedClient, accepted *packet.ConnectionAccepted) error {
	c := client.(*satellite.Client)
	c.CurrentProtocol = c.Protocol
	return nil
}

func (handler *ClientHandler) HandleConnectionDenied(client protocol.WrappedClient, denied *packet.ConnectionDenied) error {
	c := client.(*satellite.Client)
	c.Logger.Error("Connection denied")
	os.Exit(0)
	return nil
}
