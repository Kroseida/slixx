package provider

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/samsarahq/thunder/graphql"
	"gorm.io/gorm"
	"kroseida.org/slixx/pkg/model"
	_schedule "kroseida.org/slixx/pkg/schedule"
)

// JobScheduleProvider JobSchedule Provider
type JobScheduleProvider struct {
	Database *gorm.DB
}

func (provider JobScheduleProvider) Delete(id uuid.UUID) (*model.JobSchedule, error) {
	jobSchedule, err := provider.Get(id)
	if jobSchedule == nil {
		return nil, graphql.NewSafeError("jobSchedule not found")
	}
	if err != nil {
		return nil, err
	}

	result := provider.Database.Delete(&jobSchedule)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return jobSchedule, nil
}

func (provider JobScheduleProvider) Create(name string, jobId uuid.UUID, description string, kindName string, configuration string) (*model.JobSchedule, error) {
	if name == "" {
		return nil, graphql.NewSafeError("name can not be empty")
	}
	kind := _schedule.ValueOf(kindName)
	if kind == nil {
		return nil, graphql.NewSafeError("Invalid jobSchedule kind \"%s\"", kindName)
	}
	// Check if configuration is valid
	parsedConfiguration, err := kind.Parse(configuration)
	if err != nil {
		return nil, err
	}

	rawConfiguration, err := json.Marshal(parsedConfiguration)
	if err != nil {
		return nil, err
	}

	configuration = string(rawConfiguration)

	jobSchedule := model.JobSchedule{
		Id:            uuid.New(),
		JobId:         jobId,
		Name:          name,
		Description:   description,
		Kind:          kindName,
		Configuration: string(rawConfiguration),
	}

	result := provider.Database.Create(&jobSchedule)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return &jobSchedule, nil
}

func (provider JobScheduleProvider) Update(id uuid.UUID, name *string, description *string, kindName *string, configuration *string) (*model.JobSchedule, error) {
	updateJobSchedule, err := provider.Get(id)
	if updateJobSchedule == nil {
		return nil, graphql.NewSafeError("jobSchedule not found")
	}
	if err != nil {
		return nil, err
	}

	if name != nil {
		if *name == "" {
			return nil, graphql.NewSafeError("name can not be empty")
		}
		updateJobSchedule.Name = *name
	}
	if kindName != nil {
		kind := _schedule.ValueOf(*kindName)
		if kind == nil {
			return nil, graphql.NewSafeError("Invalid jobSchedule kind \"%s\"", *kindName)
		}
		updateJobSchedule.Kind = *kindName
	}
	if configuration != nil {
		kindType := _schedule.ValueOf(updateJobSchedule.Kind)

		parsedConfiguration, err := kindType.Parse(*configuration)
		if err != nil {
			return nil, err
		}

		rawConfiguration, err := json.Marshal(parsedConfiguration)
		if err != nil {
			return nil, err
		}

		updateJobSchedule.Configuration = string(rawConfiguration)
	}
	if description != nil {
		updateJobSchedule.Description = *description
	}

	result := provider.Database.Save(&updateJobSchedule)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return updateJobSchedule, nil
}

func (provider JobScheduleProvider) List() ([]*model.JobSchedule, error) {
	var jobSchedules []*model.JobSchedule
	result := provider.Database.Find(&jobSchedules)

	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return jobSchedules, nil
}

func (provider JobScheduleProvider) ListPaged(pagination *Pagination[model.JobSchedule], jobId *uuid.UUID) (*Pagination[model.JobSchedule], error) {
	var context *gorm.DB
	if jobId == nil {
		context = paginate(model.JobSchedule{}, "name", pagination, provider.Database)
	} else {
		context = paginateWithFilter(model.JobSchedule{}, "name", pagination, provider.Database, "job_id = ?", *jobId)
	}
	var schedules []model.JobSchedule
	result := context.Find(&schedules)

	if isSqlError(result.Error) {
		return nil, result.Error
	}

	pagination.Rows = schedules
	return pagination, nil
}

func (provider JobScheduleProvider) Get(id uuid.UUID) (*model.JobSchedule, error) {
	var jobSchedule *model.JobSchedule
	result := provider.Database.First(&jobSchedule, "id = ?", id)

	if isSqlError(result.Error) {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return jobSchedule, nil
}
