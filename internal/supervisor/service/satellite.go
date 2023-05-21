package service

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/internal/supervisor/datasource"
	_satellite "kroseida.org/slixx/internal/supervisor/syncnetwork"
	"kroseida.org/slixx/pkg/model"
)

func CreateSatellite(name string, description string, address string, token string) (*model.Satellite, error) {
	satellite, err := datasource.SatelliteProvider.CreateSatellite(name, description, address, token)
	if err != nil {
		return nil, err
	}

	// Create a connection to the satellite
	go _satellite.ProvideClient(*satellite)

	return satellite, err
}

func DeleteSatellite(id uuid.UUID) (*model.Satellite, error) {
	satellite, err := datasource.SatelliteProvider.DeleteSatellite(id)
	if err != nil {
		return nil, err
	}

	// Remove the connection to the satellite
	go _satellite.RemoveClient(satellite.Id)

	return satellite, err
}
