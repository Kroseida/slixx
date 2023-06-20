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
	Database          *gorm.DB
	StorageProvider   *StorageProvider
	SatelliteProvider *SatelliteProvider
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

func (provider JobProvider) CreateJob(
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

	// TODO: add this againe
	// Check if executor satellite exists
	//executorSatellite, err := provider.SatelliteProvider.GetSatellite(executorSatelliteId)
	//if executorSatellite == nil {
	//	return nil, graphql.NewSafeError("executor satellite not found")
	//}

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

func (provider JobProvider) UpdateJob(
	id uuid.UUID,
	name *string,
	description *string,
	strategyName *string,
	configuration *string,
	originStorageId *uuid.UUID,
	destinationStorageId *uuid.UUID,
	executorSatelliteId *uuid.UUID,
) (*model.Job, error) {
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
	if originStorageId != nil {
		// Check if origin storages exist
		originStorage, err := provider.StorageProvider.GetStorage(*originStorageId)
		if err != nil {
			return nil, err
		}
		if originStorage == nil {
			return nil, graphql.NewSafeError("origin storage not found")
		}
		updateJob.OriginStorageId = *originStorageId
	}
	if destinationStorageId != nil {
		// Check if destination storages exist
		destinationStorage, err := provider.StorageProvider.GetStorage(*destinationStorageId)
		if err != nil {
			return nil, err
		}
		if destinationStorage == nil {
			return nil, graphql.NewSafeError("destination storage not found")
		}
		updateJob.DestinationStorageId = *destinationStorageId
	}
	if executorSatelliteId != nil {
		// Check if executor satellite exist
		executorSatellite, err := provider.SatelliteProvider.GetSatellite(*executorSatelliteId)
		if err != nil {
			return nil, err
		}
		if executorSatellite == nil {
			return nil, graphql.NewSafeError("executor satellite not found")
		}
		updateJob.ExecutorSatelliteId = *executorSatelliteId
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

func (provider JobProvider) GetJobByStorageId(id uuid.UUID) ([]*model.Job, error) {
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
