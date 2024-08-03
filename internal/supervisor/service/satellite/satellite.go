package satellite

import (
	"github.com/google/uuid"
	"github.com/samsarahq/thunder/graphql"
	"kroseida.org/slixx/internal/supervisor/datasource"
	"kroseida.org/slixx/internal/supervisor/datasource/provider"
	"kroseida.org/slixx/internal/supervisor/slixxreactive"
	"kroseida.org/slixx/internal/supervisor/syncnetwork"
	syncnetworkClients "kroseida.org/slixx/internal/supervisor/syncnetwork/clients"
	"kroseida.org/slixx/pkg/model"
	"kroseida.org/slixx/pkg/syncnetwork/protocol"
	"time"
)

type StatefulSatellite struct {
	model.Satellite
	Connected bool
}

func StartWatchdog() {
	for {
		satellites, err := List()
		if err != nil {
			return
		}

		// Create clients for database entries that are not provided yet
		for _, satellite := range satellites {
			syncnetwork.ProvideClient(satellite.Satellite)
		}

		// Remove clients that are not in the database anymore
		satellitesMap := make(map[uuid.UUID]model.Satellite)
		for _, satellite := range satellites {
			satellitesMap[satellite.Id] = satellite.Satellite
		}
		for _, client := range syncnetworkClients.ListConnected() {
			delete(satellitesMap, client.Model.Id)
		}
		for id := range satellitesMap {
			syncnetwork.RemoveClient(id)
		}

		time.Sleep(5 * time.Minute) // Wait 5 minutes before checking again for new satellites
	}
}

func List() ([]StatefulSatellite, error) {
	satellites, err := datasource.SatelliteProvider.List()
	if err != nil {
		return nil, err
	}
	statefulSatellites := make([]StatefulSatellite, len(satellites))
	for i, satellite := range satellites {
		statefulSatellites[i] = StatefulSatellite{
			Satellite: *satellite,
			Connected: syncnetwork.GetClient(satellite.Id) != nil && syncnetwork.GetClient(satellite.Id).Client != nil &&
				syncnetwork.GetClient(satellite.Id).Client.CurrentProtocol == protocol.Supervisor,
		}
	}
	return statefulSatellites, nil
}

func Get(id uuid.UUID) (*StatefulSatellite, error) {
	satellite, err := datasource.SatelliteProvider.Get(id)
	if err != nil {
		return nil, err
	}
	return &StatefulSatellite{
		Satellite: *satellite,
		Connected: syncnetwork.GetClient(satellite.Id) != nil && syncnetwork.GetClient(satellite.Id).Client != nil &&
			syncnetwork.GetClient(satellite.Id).Client.CurrentProtocol == protocol.Supervisor,
	}, nil
}

func GetPaged(pagination *provider.Pagination[model.Satellite]) (*provider.Pagination[StatefulSatellite], error) {
	pagedData, err := datasource.SatelliteProvider.ListPaged(pagination)
	if err != nil {
		return nil, err
	}
	statefulPagedData := &provider.Pagination[StatefulSatellite]{}
	statefulPagedData.Page = pagedData.Page
	statefulPagedData.Rows = make([]StatefulSatellite, len(pagedData.Rows))
	for i, satellite := range pagedData.Rows {
		statefulPagedData.Rows[i] = StatefulSatellite{
			Satellite: satellite,
			Connected: syncnetwork.GetClient(satellite.Id) != nil && syncnetwork.GetClient(satellite.Id).Client != nil &&
				syncnetwork.GetClient(satellite.Id).Client.CurrentProtocol == protocol.Supervisor,
		}
	}
	return statefulPagedData, nil
}

func GetLogs(satelliteId uuid.UUID, pagination *provider.Pagination[model.SatelliteLogEntry]) (*provider.Pagination[model.SatelliteLogEntry], error) {
	return datasource.SatelliteProvider.GetLogs(satelliteId, pagination)
}

func Create(name string, description string, address string, token string) (*model.Satellite, error) {
	satellite, err := datasource.SatelliteProvider.Create(name, description, address, token)
	if err != nil {
		return nil, err
	}

	// Create a connection to the satellite
	go syncnetwork.ProvideClient(*satellite)

	slixxreactive.Event("satellite.create")

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

	slixxreactive.Event("satellite.update." + id.String())
	slixxreactive.Event("satellite.update.*")

	return satellite, err
}

func Delete(id uuid.UUID) (*model.Satellite, error) {
	jobs, err := datasource.JobProvider.GetByExecutorSatelliteId(id)
	if err != nil {
		return nil, err
	}
	if len(jobs) > 0 {
		return nil, graphql.NewSafeError("satellite is in use")
	}

	satellite, err := datasource.SatelliteProvider.Delete(id)
	if err != nil {
		return nil, err
	}

	// Remove the connection to the satellite
	go syncnetwork.RemoveClient(satellite.Id)

	slixxreactive.Event("satellite.update." + id.String())
	slixxreactive.Event("satellite.update.*")

	return satellite, err
}
