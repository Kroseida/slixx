package syncnetwork

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/internal/common"
	"kroseida.org/slixx/internal/supervisor/application"
	"kroseida.org/slixx/internal/supervisor/datasource"
	"kroseida.org/slixx/internal/supervisor/syncnetwork/protocol/satellite"
	"kroseida.org/slixx/internal/supervisor/syncnetwork/protocol/supervisor"
	"kroseida.org/slixx/pkg/model"
	"kroseida.org/slixx/pkg/syncnetwork"
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	"kroseida.org/slixx/pkg/syncnetwork/protocol/handshake"
	"kroseida.org/slixx/pkg/syncnetwork/protocol/supervisor/packet"
	"time"
)

var clients = make(map[uuid.UUID]*WrappedClient)
var CachedLogs []*model.SatelliteLogEntry

func Watchdog() {
	for {
		satellites, err := datasource.SatelliteProvider.GetSatellites()
		if err != nil {
			return
		}

		// Create clients for database entries that are not provided yet
		for _, satellite := range satellites {
			ProvideClient(*satellite)
		}

		// Remove clients that are not in the database anymore
		satellitesMap := make(map[uuid.UUID]*model.Satellite)
		for _, satellite := range satellites {
			satellitesMap[satellite.Id] = satellite
		}
		for _, client := range clients {
			delete(satellitesMap, client.Model.Id)
		}
		for id := range satellitesMap {
			RemoveClient(id)
		}

		time.Sleep(5 * time.Minute) // Wait 5 minutes before checking again for new satellites
	}
}

func RemoveClient(id uuid.UUID) {
	client := clients[id]
	if client == nil {
		return
	}
	client.Client.Close()
	delete(clients, id)
}

func GetClient(id uuid.UUID) *WrappedClient {
	return clients[id]
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
			SyncStorages()
			SyncJobs()
		},
		Version: common.CurrentVersion,
	}
	// TODO: make timeout configurable in the database or in the configuration file - not sure yet
	go client.Dial(5*time.Second, 5*time.Second)

	clients[configuration.Id] = &WrappedClient{
		Model:  configuration,
		Client: &client,
	}
}

func ApplyUpdates(configuration model.Satellite) {
	client := clients[configuration.Id]
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
	_, bool := clients[id]
	return bool
}
