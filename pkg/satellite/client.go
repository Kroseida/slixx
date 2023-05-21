package satellite

import (
	"bufio"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm/utils"
	"kroseida.org/slixx/internal/supervisor/application"
	"kroseida.org/slixx/pkg/satellite/protocol"
	"kroseida.org/slixx/pkg/satellite/protocol/handshake/packet"
	"net"
	"strconv"
	"time"
)

type Client struct {
	Address         string
	connection      net.Conn
	Closed          bool
	Logger          *zap.SugaredLogger
	Token           string
	Protocol        string
	CurrentProtocol string
	Handler         map[string]protocol.Handler
	Reader          *bufio.Reader
	Writer          *bufio.Writer
}

func (client *Client) Close() {
	client.Closed = true
	if client.connection == nil {
		return
	}
	err := client.connection.Close()
	if err != nil {
		client.Logger.Error("Error while closing sync network connection", err)
	}
}

func (client *Client) Dial(timeout time.Duration, reconnectAfter time.Duration) {
	connection, err := net.DialTimeout("tcp", client.Address, timeout)
	if err != nil {
		client.Logger.Error(
			"Failed to connect to satellite sync network ("+client.Address+") retrying in "+strconv.Itoa(int(reconnectAfter.Milliseconds()))+"ms: ",
			err,
		)
		time.Sleep(reconnectAfter)
		if client.Closed {
			return
		}
		client.Dial(timeout, reconnectAfter)
		return
	}
	client.connection = connection
	client.CurrentProtocol = protocol.HandshakeProtocol
	client.Reader = bufio.NewReader(connection)
	client.Writer = bufio.NewWriter(connection)

	// Send handshake packet
	err = client.Send(&packet.Handshake{
		Token:          client.Token,
		TargetProtocol: client.Protocol,
		Version:        "1.0.0", // TODO: Get version from somewhere else
	})
	if err != nil {
		application.Logger.Error("Failed to send handshake packet to satellite sync network (" + client.Address + "): " + err.Error())
		client.Close()
	}

	err = client.listen()

	if client.Closed {
		return
	}
	if err != nil {
		client.Logger.Error(
			"Lost connection to satellite sync network ("+client.Address+") retrying in "+strconv.Itoa(int(reconnectAfter.Milliseconds()))+"ms: ",
			err,
		)
	}
	time.Sleep(reconnectAfter)
	client.Dial(timeout, reconnectAfter)
}

func (client *Client) Send(packet protocol.Packet) error {
	if !utils.Contains(packet.Protocol(), client.CurrentProtocol) {
		return errors.New("Packet with id " + strconv.Itoa(int(packet.PacketId())) + " is not supported by the current protocol (" + client.CurrentProtocol + ")")
	}
	return protocol.SendPacket(*client.Writer, packet)
}

func (client *Client) listen() error {
	client.Logger.Info("Connected to satellite sync network (" + client.Address + ")")
	for !client.Closed {
		packet, err := protocol.ReadPacket(client.Reader, PACKETS)
		if err != nil {
			return err
		}
		handler, ok := client.Handler[client.Protocol]
		if !ok {
			client.Logger.Error("No handler for protocol " + client.Protocol)
			continue
		}
		err = handler.Handle(client, packet)
		if err != nil {
			client.Logger.Error("Error while handling packet: ", err)
		}
	}
	return nil
}
