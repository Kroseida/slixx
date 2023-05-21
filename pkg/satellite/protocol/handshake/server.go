package handshake

import (
	"kroseida.org/slixx/pkg/satellite"
	"kroseida.org/slixx/pkg/satellite/protocol"
	"kroseida.org/slixx/pkg/satellite/protocol/handshake/packet"
)

type ServerHandler struct {
	Token string
}

func (handler *ServerHandler) Handle(client protocol.WrappedClient, p protocol.Packet) error {
	if p.PacketId() == (&packet.Handshake{}).PacketId() {
		return handler.HandleHandshake(client, p.(*packet.Handshake))
	}

	return nil
}

func (handler *ServerHandler) HandleHandshake(client protocol.WrappedClient, handshake *packet.Handshake) error {
	c := client.(*satellite.ConnectedClient)

	if handshake.Token != handler.Token {
		return c.Send(&packet.ConnectionDenied{})
	}

	c.Id = &handshake.Id
	c.Protocol = handshake.TargetProtocol
	return c.Send(&packet.ConnectionAccepted{})
}
