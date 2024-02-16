package job_schedule

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/internal/supervisor/datasource"
	"kroseida.org/slixx/internal/supervisor/datasource/provider"
	"kroseida.org/slixx/internal/supervisor/syncnetwork/action"
	"kroseida.org/slixx/pkg/model"
)

func Get(id uuid.UUID) (*model.JobSchedule, error) {
	return datasource.JobScheduleProvider.Get(id)
}

func GetPaged(pagination *provider.Pagination[model.JobSchedule], jobId *uuid.UUID) (*provider.Pagination[model.JobSchedule], error) {
	return datasource.JobScheduleProvider.ListPaged(pagination, jobId)
}

func Create(name string, jobId uuid.UUID, description string, kindName string, configuration string) (*model.JobSchedule, error) {
	jobSchedule, err := datasource.JobScheduleProvider.Create(name, jobId, description, kindName, configuration)
	if err != nil {
		return nil, err
	}

	go action.SyncJobSchedules(nil)

	return jobSchedule, nil
}

func Update(id uuid.UUID, name *string, description *string, kindName *string, configuration *string) (*model.JobSchedule, error) {
	jobSchedule, err := datasource.JobScheduleProvider.Update(
		id,
		name,
		description,
		kindName,
		configuration,
	)

	if err != nil {
		return nil, err
	}

	go action.SyncJobSchedules(nil)

	return jobSchedule, nil
}

func Delete(id uuid.UUID) (*model.JobSchedule, error) {
	jobSchedule, err := datasource.JobScheduleProvider.Delete(id)
	if err != nil {
		return nil, err
	}

	go action.SyncJobSchedules(nil)

	return jobSchedule, nil
}
