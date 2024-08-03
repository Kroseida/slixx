package handshake

import (
	"kroseida.org/slixx/pkg/syncnetwork"
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	"kroseida.org/slixx/pkg/syncnetwork/protocol/handshake/packet"
)

type ClientHandler struct {
}

func (handler *ClientHandler) Handle(client protocol.WrappedClient, p protocol.Packet) error {
	if p.PacketId() == (&packet.ConnectionAccepted{}).PacketId() {
		return handler.HandleConnectionAccepted(client)
	}
	if p.PacketId() == (&packet.ConnectionDenied{}).PacketId() {
		return handler.HandleConnectionDenied(client)
	}

	return nil
}

func (handler *ClientHandler) HandleConnectionAccepted(client protocol.WrappedClient) error {
	c := client.(*syncnetwork.Client)
	c.CurrentProtocol = c.Protocol

	if c.AfterProtocolSelection != nil {
		c.Connected = true
		c.AfterProtocolSelection(client)
	}

	return nil
}

func (handler *ClientHandler) HandleConnectionDenied(client protocol.WrappedClient) error {
	c := client.(*syncnetwork.Client)
	c.LogListener("error", "Connection denied to remote ("+c.Connection.RemoteAddr().String()+")")
	c.Connection.Close() // Force to reconnect
	return nil
}
