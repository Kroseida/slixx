package syncnetwork

import (
	"bufio"
	"errors"
	gormUtils "gorm.io/gorm/utils"
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	"kroseida.org/slixx/pkg/syncnetwork/protocol/handshake/packet"
	"net"
	"strconv"
	"time"
)

type Client struct {
	Address                string
	Connection             net.Conn
	Closed                 bool
	Token                  string
	Protocol               string
	CurrentProtocol        string
	Handler                map[string]protocol.Handler
	Reader                 *bufio.Reader
	Writer                 *bufio.Writer
	AfterProtocolSelection func(protocol.WrappedClient)
	Version                string
	LogListener            func(string, ...interface{})
}

func (client *Client) Close() {
	client.Closed = true
	if client.Connection == nil {
		return
	}
	err := client.Connection.Close()
	if err != nil {
		client.LogListener("error", "Error while closing sync network connection", err)
	}
}

func (client *Client) Dial(timeout time.Duration, reconnectAfter time.Duration) {
	for true {
		connection, err := net.DialTimeout("tcp", client.Address, timeout)
		if err != nil {
			client.LogListener(
				"warn",
				"Failed to connect to satellite sync network ("+client.Address+") retrying in "+strconv.Itoa(int(reconnectAfter.Milliseconds()))+"ms: ",
				err,
			)
			time.Sleep(reconnectAfter)
			if client.Closed {
				return
			}
			continue
		}
		client.Connection = connection
		client.CurrentProtocol = protocol.Handshake
		client.Reader = bufio.NewReader(connection)
		client.Writer = bufio.NewWriter(connection)

		// Send handshake packet
		err = client.Send(&packet.Handshake{
			Token:          client.Token,
			TargetProtocol: client.Protocol,
			Version:        client.Version,
		})
		if err != nil {
			client.LogListener("error", "Failed to send handshake packet *to satellite sync network ("+client.Address+"): "+err.Error())
			client.Close()
		}

		err = client.listen()

		client.CurrentProtocol = ""
		client.Connection = nil
		client.Reader = nil
		client.Writer = nil
		if client.Closed {
			return
		}
		if err != nil {
			client.LogListener(
				"error",
				"Lost connection to satellite sync network ("+client.Address+") retrying in "+strconv.Itoa(int(reconnectAfter.Milliseconds()))+"ms: ",
				err,
			)
		}
		time.Sleep(reconnectAfter)
		continue
	}
}

func (client *Client) Send(packet protocol.Packet) error {
	if !gormUtils.Contains(packet.Protocol(), client.CurrentProtocol) {
		return errors.New("Packet with id " + strconv.Itoa(int(packet.PacketId())) + " is not supported by the current protocol (" + client.CurrentProtocol + ")")
	}
	return protocol.SendPacket(*client.Writer, packet)
}

func (client *Client) listen() error {
	client.LogListener("info", "Connected to sync network ("+client.Address+")")
	for !client.Closed {
		packet, err := protocol.ReadPacket(client.Reader, PACKETS)
		if err != nil {
			return err
		}
		handler, ok := client.Handler[client.CurrentProtocol]
		if !ok {
			client.LogListener("warning", "No handler for protocol "+client.CurrentProtocol)
			continue
		}
		err = handler.Handle(client, packet)
		if err != nil {
			client.LogListener("error", "Error while handling packet *with id ("+strconv.Itoa(int(packet.PacketId()))+"): ", err)
		}
	}
	client.CurrentProtocol = ""
	client.Connection = nil
	client.Reader = nil
	client.Writer = nil

	return nil
}
