package satellitelogservice

import (
	"kroseida.org/slixx/internal/supervisor/datasource"
	"kroseida.org/slixx/pkg/model"
)

func Create(logs []*model.SatelliteLogEntry) error {
	return datasource.SatelliteProvider.CreateSatelliteLogs(logs)
}
