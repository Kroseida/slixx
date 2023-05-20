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
	Database        *gorm.DB
	StorageProvider *StorageProvider
}

func (provider JobProvider) DeleteJob(id uuid.UUID) (*model.Job, error) {
	job, err := provider.GetJob(id)
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

func (provider JobProvider) CreateJob(name string, description string, strategyName string, configuration string, originStorageId uuid.UUID, destinationStorageId uuid.UUID) (*model.Job, error) {
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

	// Check if origin storages exist
	originStorage, err := provider.StorageProvider.GetStorage(originStorageId)
	if err != nil {
		return nil, err
	}
	if originStorage == nil {
		return nil, graphql.NewSafeError("origin storage not found")
	}

	// Check if destination storages exist
	destinationStorage, err := provider.StorageProvider.GetStorage(destinationStorageId)
	if err != nil {
		return nil, err
	}
	if destinationStorage == nil {
		return nil, graphql.NewSafeError("destination storage not found")
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
	}

	result := provider.Database.Create(&job)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return &job, nil
}

func (provider JobProvider) UpdateJob(id uuid.UUID, name *string, description *string, strategyName *string, configuration *string, originStorage *uuid.UUID, destinationStorage *uuid.UUID) (*model.Job, error) {
	updateJob, err := provider.GetJob(id)
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
	if originStorage != nil {
		updateJob.OriginStorageId = *originStorage
	}
	if destinationStorage != nil {
		updateJob.DestinationStorageId = *destinationStorage
	}

	result := provider.Database.Save(&updateJob)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return updateJob, nil
}

func (provider JobProvider) GetJobs() ([]*model.Job, error) {
	var jobs []*model.Job
	result := provider.Database.Find(&jobs)

	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return jobs, nil
}

func (provider JobProvider) GetJobsPaged(pagination *Pagination[model.Job]) (*Pagination[model.Job], error) {
	context := paginate(model.Job{}, "name", pagination, provider.Database)

	var jobs []model.Job
	result := context.Find(&jobs)

	if isSqlError(result.Error) {
		return nil, result.Error
	}

	pagination.Rows = jobs
	return pagination, nil
}

func (provider JobProvider) GetJob(id uuid.UUID) (*model.Job, error) {
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
