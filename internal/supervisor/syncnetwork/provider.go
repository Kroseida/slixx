package syncnetwork

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/internal/satellite/syncnetwork/protocol/supervisor"
	"kroseida.org/slixx/internal/supervisor/application"
	"kroseida.org/slixx/internal/supervisor/datasource"
	"kroseida.org/slixx/pkg/model"
	"kroseida.org/slixx/pkg/syncnetwork"
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	"kroseida.org/slixx/pkg/syncnetwork/protocol/handshake"
	"time"
)

var clients = make(map[uuid.UUID]*WrappedClient)

func Watchdog() {
	for true {
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
		Protocol: protocol.SupervisorProtocol,
		Handler: map[string]protocol.Handler{
			protocol.HandshakeProtocol:  &handshake.ClientHandler{},
			protocol.SupervisorProtocol: &supervisor.Handler{},
			protocol.SatelliteProtocol:  &supervisor.Handler{},
		},
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
