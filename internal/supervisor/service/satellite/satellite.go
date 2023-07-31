package satellite

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/internal/supervisor/datasource"
	"kroseida.org/slixx/internal/supervisor/datasource/provider"
	"kroseida.org/slixx/internal/supervisor/syncnetwork"
	syncnetworkClients "kroseida.org/slixx/internal/supervisor/syncnetwork/clients"
	"kroseida.org/slixx/pkg/model"
	"time"
)

func StartWatchdog() {
	for {
		satellites, err := List()
		if err != nil {
			return
		}

		// Create clients for database entries that are not provided yet
		for _, satellite := range satellites {
			syncnetwork.ProvideClient(*satellite)
		}

		// Remove clients that are not in the database anymore
		satellitesMap := make(map[uuid.UUID]*model.Satellite)
		for _, satellite := range satellites {
			satellitesMap[satellite.Id] = satellite
		}
		for _, client := range syncnetworkClients.List {
			delete(satellitesMap, client.Model.Id)
		}
		for id := range satellitesMap {
			syncnetwork.RemoveClient(id)
		}

		time.Sleep(5 * time.Minute) // Wait 5 minutes before checking again for new satellites
	}
}

func List() ([]*model.Satellite, error) {
	return datasource.SatelliteProvider.List()
}

func Get(id uuid.UUID) (*model.Satellite, error) {
	return datasource.SatelliteProvider.Get(id)
}

func GetPaged(pagination *provider.Pagination[model.Satellite]) (*provider.Pagination[model.Satellite], error) {
	return datasource.SatelliteProvider.ListPaged(pagination)
}

func Create(name string, description string, address string, token string) (*model.Satellite, error) {
	satellite, err := datasource.SatelliteProvider.Create(name, description, address, token)
	if err != nil {
		return nil, err
	}

	// Create a connection to the satellite
	go syncnetwork.ProvideClient(*satellite)

	return satellite, err
}

func Update(id uuid.UUID, name *string, description *string, address *string, token *string) (*model.Satellite, error) {
	satellite, err := datasource.SatelliteProvider.Update(id, name, description, address, token)
	if err != nil {
		return nil, err
	}

	go func() {
		syncnetwork.ApplyUpdates(*satellite)
		client := syncnetwork.GetClient(id)
		if client != nil && client.Client != nil && client.Client.Connection != nil {
			client.Client.Connection.Close() // force a reconnect
		}
	}()

	return satellite, err
}

func Delete(id uuid.UUID) (*model.Satellite, error) {
	satellite, err := datasource.SatelliteProvider.Delete(id)
	if err != nil {
		return nil, err
	}

	// Remove the connection to the satellite
	go syncnetwork.RemoveClient(satellite.Id)

	return satellite, err
}
