package handshake

import (
	"kroseida.org/slixx/pkg/syncnetwork"
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	"kroseida.org/slixx/pkg/syncnetwork/protocol/handshake/packet"
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
	c := client.(*syncnetwork.ConnectedClient)

	if handshake.Token != handler.Token {
		c.Server.Logger.Warn("Connection of client(" + (*c.Connection).RemoteAddr().String() + ") denied: Invalid token")
		return c.Send(&packet.ConnectionDenied{})
	}
	if handshake.Version != c.Server.Version {
		c.Server.Logger.Warn("Connection of client(" + (*c.Connection).RemoteAddr().String() + ") denied: Invalid version")
		c.Server.Logger.Warn("Client version: " + handshake.Version + " | Server version: " + c.Server.Version)
		// TODO: maybe we will add a feature to auto update the satellite via the supervisor
		return c.Send(&packet.ConnectionDenied{})
	}

	c.Id = &handshake.Id
	c.Protocol = handshake.TargetProtocol
	c.Server.Logger.Info("Connection of client(" + (*c.Connection).RemoteAddr().String() + ") accepted as " + c.Protocol)

	err := c.Send(&packet.ConnectionAccepted{})
	if err != nil {
		return err
	}

	if c.Server.AfterProtocolSelection != nil {
		c.Server.AfterProtocolSelection(client)
	}
	return nil
}
