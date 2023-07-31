package syncnetwork

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/internal/common"
	"kroseida.org/slixx/internal/supervisor/application"
	"kroseida.org/slixx/internal/supervisor/syncnetwork/action"
	"kroseida.org/slixx/internal/supervisor/syncnetwork/clients"
	supervisorProtocol "kroseida.org/slixx/internal/supervisor/syncnetwork/protocol"
	"kroseida.org/slixx/internal/supervisor/syncnetwork/protocol/satellite"
	"kroseida.org/slixx/internal/supervisor/syncnetwork/protocol/supervisor"
	"kroseida.org/slixx/pkg/model"
	"kroseida.org/slixx/pkg/syncnetwork"
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	"kroseida.org/slixx/pkg/syncnetwork/protocol/handshake"
	"kroseida.org/slixx/pkg/syncnetwork/protocol/supervisor/packet"
	"time"
)

func RemoveClient(id uuid.UUID) {
	client := clients.List[id]
	if client == nil {
		return
	}
	client.Client.Close()
	delete(clients.List, id)
}

func GetClient(id uuid.UUID) *supervisorProtocol.WrappedClient {
	return clients.List[id]
}

func ProvideClient(configuration model.Satellite) {
	if IsProvided(configuration.Id) {
		ApplyUpdates(configuration)
		return
	}
	client := syncnetwork.Client{
		Address:  configuration.Address,
		Token:    configuration.Token,
		Closed:   false,
		Logger:   application.Logger,
		Protocol: protocol.Supervisor,
		Handler: map[string]protocol.Handler{
			protocol.Handshake:  &handshake.ClientHandler{},
			protocol.Supervisor: &supervisor.Handler{},
			protocol.Satellite:  &satellite.Handler{},
		},
		AfterProtocolSelection: func(client protocol.WrappedClient) {
			client.Send(&packet.ApplySupervisor{
				Id: configuration.Id,
			})
			action.SyncStorages()
			action.SyncJobs()
		},
		Version: common.CurrentVersion,
	}
	// TODO: make timeout configurable in the database or in the configuration file - not sure yet
	go client.Dial(5*time.Second, 5*time.Second)

	clients.List[configuration.Id] = &supervisorProtocol.WrappedClient{
		Model:  configuration,
		Client: &client,
	}
}

func ApplyUpdates(configuration model.Satellite) {
	client := clients.List[configuration.Id]
	if client == nil {
		return
	}
	if client.Model.Address != configuration.Address {
		client.Model.Address = configuration.Address
		client.Client.Address = configuration.Address
	}
	if client.Model.Token != configuration.Token {
		client.Model.Token = configuration.Token
		client.Client.Token = configuration.Token
	}
}

func IsProvided(id uuid.UUID) bool {
	_, bool := clients.List[id]
	return bool
}
