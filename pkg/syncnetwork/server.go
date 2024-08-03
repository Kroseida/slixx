package syncnetwork

import (
	"bufio"
	"errors"
	"github.com/google/uuid"
	gormUtils "gorm.io/gorm/utils"
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	"kroseida.org/slixx/pkg/utils"
	"net"
	"strconv"
)

type Server struct {
	Id                     *uuid.UUID
	BindAddress            string
	Token                  string
	listener               net.Listener
	Handler                map[string]protocol.Handler
	closed                 bool
	Logger                 utils.Logger
	Version                string
	ActiveConnection       []*ConnectedClient
	AfterProtocolSelection func(protocol.WrappedClient)
}

type ConnectedClient struct {
	Connection *net.Conn
	Protocol   string
	Id         *string
	Reader     *bufio.Reader
	Writer     *bufio.Writer
	Server     *Server
	Connected  bool
}

func (client *ConnectedClient) IsConnected() bool {
	return client.Connected
}

func (client *ConnectedClient) Send(packet protocol.Packet) error {
	if !gormUtils.Contains(packet.Protocol(), client.Protocol) {
		return errors.New("Packet with id " + strconv.Itoa(int(packet.PacketId())) + " is not supported by the current protocol (" + client.Protocol + ")")
	}
	return protocol.SendPacket(*client.Writer, packet)
}

func (server *Server) Close(client *ConnectedClient) error {
	for i, c := range server.ActiveConnection {
		if c.Connection == client.Connection {
			server.ActiveConnection = append(server.ActiveConnection[:i], server.ActiveConnection[i+1:]...)
			break
		}
	}
	return (*client.Connection).Close()
}

func (server *Server) BroadcastSatellites(packet protocol.Packet) error {
	for _, client := range server.ActiveConnection {
		if client.Protocol != protocol.Satellite {
			continue
		}
		err := client.Send(packet)
		if err != nil {
			server.Logger.Error("Error while broadcasting packet *to client", err)
			continue // We don't want to stop the broadcast if one client fails
		}
	}
	return nil
}

func (server *Server) Broadcast(packet protocol.Packet) error {
	for _, client := range server.ActiveConnection {
		err := client.Send(packet)
		if err != nil {
			server.Logger.Error("Error while broadcasting packet *to client", err)
			continue // We don't want to stop the broadcast if one client fails
		}
	}
	return nil
}

func (server *Server) Listen() error {
	listener, err := net.Listen("tcp", server.BindAddress)
	if err != nil {
		return err
	}
	server.listener = listener
	for !server.closed {
		connection, err := server.listener.Accept()
		if err != nil {
			return err
		}

		server.Logger.Info("New connection from: ", connection.RemoteAddr().String())
		client := ConnectedClient{
			Connection: &connection,
			Protocol:   protocol.Handshake,
			Id:         nil,
			Reader:     bufio.NewReader(connection),
			Writer:     bufio.NewWriter(connection),
			Server:     server,
		}
		server.ActiveConnection = append(server.ActiveConnection, &client)

		go server.handleConnection(&client)
	}
	return nil
}

func (server *Server) handleConnection(client *ConnectedClient) {
	for true {
		packet, err := protocol.ReadPacket(client.Reader, PACKETS)
		if err != nil {
			// Close connection if we can't read a packet *so we don't get stuck in a loop and force a reconnect
			server.Close(client)
			return
		}
		handler, ok := server.Handler[client.Protocol]
		if !ok {
			server.Logger.Error("No handler for protocol " + client.Protocol)
			continue
		}
		err = handler.Handle(client, packet)
		if err != nil {
			server.Logger.Error("Error while handling packet: ", err)
		}
	}
}
