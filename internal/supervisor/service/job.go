package service

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/internal/supervisor/datasource"
	"kroseida.org/slixx/internal/supervisor/syncnetwork"
	"kroseida.org/slixx/pkg/model"
)

func CreateJob(
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

	syncnetwork.SyncJobs()

	return job, err
}

func UpdateJob(
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

	syncnetwork.SyncJobs()

	return job, err
}

func DeleteJob(id uuid.UUID) (*model.Job, error) {
	job, err := datasource.JobProvider.DeleteJob(id)
	if err != nil {
		return nil, err
	}
	
	syncnetwork.SyncJobs()

	return job, err
}
