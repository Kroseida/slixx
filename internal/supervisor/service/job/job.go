package job

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/internal/supervisor/datasource"
	"kroseida.org/slixx/internal/supervisor/datasource/provider"
	"kroseida.org/slixx/internal/supervisor/syncnetwork/action"
	"kroseida.org/slixx/pkg/model"
)

func Get(id uuid.UUID) (*model.Job, error) {
	return datasource.JobProvider.GetJob(id)
}

func GetPaged(pagination *provider.Pagination[model.Job]) (*provider.Pagination[model.Job], error) {
	return datasource.JobProvider.GetJobsPaged(pagination)
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
	job, err := datasource.JobProvider.CreateJob(
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

	action.SyncJobs()

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
	job, err := datasource.JobProvider.UpdateJob(
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

	action.SyncJobs()

	return job, err
}

func Delete(id uuid.UUID) (*model.Job, error) {
	job, err := datasource.JobProvider.DeleteJob(id)
	if err != nil {
		return nil, err
	}

	action.SyncJobs()

	return job, err
}
