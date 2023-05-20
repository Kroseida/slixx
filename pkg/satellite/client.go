package satellite

import (
	"go.uber.org/zap"
	"net"
	"time"
)

type Client struct {
	Address    string
	connection net.Conn
	Closed     bool
	Logger     *zap.SugaredLogger
}

func (client *Client) Close() {
	client.Closed = true
	err := client.connection.Close()
	if err != nil {
		client.Logger.Error("Error while closing connection", err)
	}
}

func (client *Client) Dial(timeout time.Duration, reconnectAfter time.Duration) {
	connection, err := net.DialTimeout("tcp", client.Address, timeout)
	if err != nil {
		client.Logger.Error("Failed to connect to satellite ("+client.Address+") retrying in "+string(reconnectAfter.Milliseconds())+"ms", err)
		time.Sleep(reconnectAfter)
		client.Dial(timeout, reconnectAfter)
		return
	}
	client.connection = connection

	err = client.listen()
	if client.Closed {
		return
	}
	if err != nil {
		client.Logger.Error("Lost connection to satellite ("+client.Address+") retrying in "+string(reconnectAfter.Milliseconds())+"ms", err)
	}
	client.Dial(timeout, reconnectAfter)
}

func (client *Client) listen() error {
	for !client.Closed {
		data := make([]byte, 4096)
		_, err := client.connection.Read(data)
		if err != nil {
			return err
		}
	}
	return nil
}
