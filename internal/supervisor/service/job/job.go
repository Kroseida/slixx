package job

import (
	"github.com/google/uuid"
	"github.com/samsarahq/thunder/graphql"
	"kroseida.org/slixx/internal/supervisor/datasource"
	"kroseida.org/slixx/internal/supervisor/datasource/provider"
	"kroseida.org/slixx/internal/supervisor/syncnetwork/action"
	"kroseida.org/slixx/pkg/model"
)

func Get(id uuid.UUID) (*model.Job, error) {
	return datasource.JobProvider.Get(id)
}

func GetPaged(pagination *provider.Pagination[model.Job]) (*provider.Pagination[model.Job], error) {
	return datasource.JobProvider.ListPaged(pagination)
}

func Create(
	name string,
	description string,
	strategyName string,
	configuration string,
	originStorageId uuid.UUID,
	destinationStorageId uuid.UUID,
	executorSatelliteId uuid.UUID,
) (*model.Job, error) {
	// Check if origin storages exist
	originStorage, err := datasource.StorageProvider.Get(originStorageId)
	if err != nil {
		return nil, err
	}
	if originStorage == nil {
		return nil, graphql.NewSafeError("origin storage not found")
	}

	// Check if destination storages exist
	destinationStorage, err := datasource.StorageProvider.Get(destinationStorageId)
	if err != nil {
		return nil, err
	}
	if destinationStorage == nil {
		return nil, graphql.NewSafeError("destination storage not found")
	}

	// Check if executor satellite exists
	executorSatellite, err := datasource.SatelliteProvider.Get(executorSatelliteId)
	if executorSatellite == nil {
		return nil, graphql.NewSafeError("executor satellite not found")
	}

	job, err := datasource.JobProvider.Create(
		name,
		description,
		strategyName,
		configuration,
		originStorageId,
		destinationStorageId,
		executorSatelliteId,
	)
	if err != nil {
		return nil, err
	}

	action.SyncJobs(nil)

	return job, err
}

func Update(
	id uuid.UUID,
	name *string,
	description *string,
	strategyName *string,
	configuration *string,
	originStorageId *uuid.UUID,
	destinationStorageId *uuid.UUID,
	executorSatelliteId *uuid.UUID,
) (*model.Job, error) {
	if originStorageId != nil {
		// Check if origin storages exist
		originStorage, err := datasource.StorageProvider.Get(*originStorageId)
		if err != nil {
			return nil, err
		}
		if originStorage == nil {
			return nil, graphql.NewSafeError("origin storage not found")
		}
	}
	if destinationStorageId != nil {
		// Check if destination storages exist
		destinationStorage, err := datasource.StorageProvider.Get(*destinationStorageId)
		if err != nil {
			return nil, err
		}
		if destinationStorage == nil {
			return nil, graphql.NewSafeError("destination storage not found")
		}
	}
	if executorSatelliteId != nil {
		// Check if executor satellite exist
		executorSatellite, err := datasource.SatelliteProvider.Get(*executorSatelliteId)
		if err != nil {
			return nil, err
		}
		if executorSatellite == nil {
			return nil, graphql.NewSafeError("executor satellite not found")
		}
	}

	job, err := datasource.JobProvider.Update(
		id,
		name,
		description,
		strategyName,
		configuration,
		originStorageId,
		destinationStorageId,
		executorSatelliteId,
	)
	if err != nil {
		return nil, err
	}

	action.SyncJobs(nil)

	return job, err
}

func Delete(id uuid.UUID) (*model.Job, error) {
	job, err := datasource.JobProvider.Delete(id)
	if err != nil {
		return nil, err
	}

	action.SyncJobs(nil)

	return job, err
}
