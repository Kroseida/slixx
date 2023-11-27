package provider

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/samsarahq/thunder/graphql"
	"gorm.io/gorm"
	"kroseida.org/slixx/pkg/model"
	_strategy "kroseida.org/slixx/pkg/strategy"
)

// JobProvider Job Provider
type JobProvider struct {
	Database *gorm.DB
}

func (provider JobProvider) Delete(id uuid.UUID) (*model.Job, error) {
	job, err := provider.Get(id)
	if job == nil {
		return nil, graphql.NewSafeError("job not found")
	}
	if err != nil {
		return nil, err
	}

	result := provider.Database.Delete(&job)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return job, nil
}

func (provider JobProvider) Create(
	name string,
	description string,
	strategyName string,
	configuration string,
	originStorageId uuid.UUID,
	destinationStorageId uuid.UUID,
	executorSatelliteId uuid.UUID,
) (*model.Job, error) {
	if name == "" {
		return nil, graphql.NewSafeError("name can not be empty")
	}
	strategy := _strategy.ValueOf(strategyName)
	if strategy == nil {
		return nil, graphql.NewSafeError("Invalid strategy \"%s\"", strategyName)
	}
	parsedConfiguration, err := strategy.Parse(configuration)
	if err != nil {
		return nil, err
	}

	rawConfiguration, err := json.Marshal(parsedConfiguration)
	if err != nil {
		return nil, err
	}

	configuration = string(rawConfiguration)

	job := model.Job{
		Id:                   uuid.New(),
		Name:                 name,
		Description:          description,
		Strategy:             strategyName,
		Configuration:        configuration,
		OriginStorageId:      originStorageId,
		DestinationStorageId: destinationStorageId,
		ExecutorSatelliteId:  executorSatelliteId,
	}

	result := provider.Database.Create(&job)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return &job, nil
}

func (provider JobProvider) Update(
	id uuid.UUID,
	name *string,
	description *string,
	strategyName *string,
	configuration *string,
	originStorageId *uuid.UUID,
	destinationStorageId *uuid.UUID,
	executorSatelliteId *uuid.UUID,
) (*model.Job, error) {
	updateJob, err := provider.Get(id)
	if updateJob == nil {
		return nil, graphql.NewSafeError("job not found")
	}
	if err != nil {
		return nil, err
	}

	if name != nil {
		if *name == "" {
			return nil, graphql.NewSafeError("name can not be empty")
		}
		updateJob.Name = *name
	}
	if strategyName != nil {
		strategy := _strategy.ValueOf(*strategyName)
		if strategy == nil {
			return nil, graphql.NewSafeError("Invalid strategy \"%s\"", *strategyName)
		}
		updateJob.Strategy = *strategyName
	}
	if configuration != nil {
		strategy := _strategy.ValueOf(updateJob.Strategy)

		parsedConfiguration, err := strategy.Parse(*configuration)
		if err != nil {
			return nil, err
		}

		rawConfiguration, err := json.Marshal(parsedConfiguration)
		if err != nil {
			return nil, err
		}

		updateJob.Configuration = string(rawConfiguration)
	}
	if description != nil {
		updateJob.Description = *description
	}
	if originStorageId != nil {
		updateJob.OriginStorageId = *originStorageId
	}
	if destinationStorageId != nil {
		updateJob.DestinationStorageId = *destinationStorageId
	}
	if executorSatelliteId != nil {
		updateJob.ExecutorSatelliteId = *executorSatelliteId
	}

	result := provider.Database.Save(&updateJob)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return updateJob, nil
}

func (provider JobProvider) List() ([]*model.Job, error) {
	var jobs []*model.Job
	result := provider.Database.Find(&jobs)

	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return jobs, nil
}

func (provider JobProvider) ListPaged(pagination *Pagination[model.Job]) (*Pagination[model.Job], error) {
	context := paginate(model.Job{}, "name", pagination, provider.Database)

	var jobs []model.Job
	result := context.Find(&jobs)

	if isSqlError(result.Error) {
		return nil, result.Error
	}

	pagination.Rows = jobs
	return pagination, nil
}

func (provider JobProvider) Get(id uuid.UUID) (*model.Job, error) {
	var job *model.Job
	result := provider.Database.First(&job, "id = ?", id)

	if isSqlError(result.Error) {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return job, nil
}

func (provider JobProvider) GetByStorageId(id uuid.UUID) ([]*model.Job, error) {
	var job []*model.Job
	result := provider.Database.Find(&job, "origin_storage_id = ? OR destination_storage_id = ?", id, id)
	if isSqlError(result.Error) {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return job, nil
}

func (provider JobProvider) GetByExecutorSatelliteId(id uuid.UUID) ([]*model.Job, error) {
	var job []*model.Job
	result := provider.Database.Find(&job, "executor_satellite_id = ?", id)
	if isSqlError(result.Error) {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return job, nil
}
