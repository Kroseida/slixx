package satelliteservice

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/internal/supervisor/datasource"
	"kroseida.org/slixx/internal/supervisor/syncnetwork"
	"kroseida.org/slixx/pkg/model"
)

func Create(name string, description string, address string, token string) (*model.Satellite, error) {
	satellite, err := datasource.SatelliteProvider.CreateSatellite(name, description, address, token)
	if err != nil {
		return nil, err
	}

	// Create a connection to the satellite
	go syncnetwork.ProvideClient(*satellite)

	return satellite, err
}

func Update(id uuid.UUID, name *string, description *string, address *string, token *string) (*model.Satellite, error) {
	satellite, err := datasource.SatelliteProvider.UpdateSatellite(id, name, description, address, token)
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
	satellite, err := datasource.SatelliteProvider.DeleteSatellite(id)
	if err != nil {
		return nil, err
	}

	// Remove the connection to the satellite
	go syncnetwork.RemoveClient(satellite.Id)

	return satellite, err
}
