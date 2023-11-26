package provider

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"kroseida.org/slixx/pkg/model"
	"time"
)

// ExecutionProvider Execution Provider
type ExecutionProvider struct {
	Database *gorm.DB
}

func (provider ExecutionProvider) Get(id uuid.UUID) (*model.Execution, error) {
	var execution *model.Execution
	result := provider.Database.First(&execution, "id = ?", id)

	if isSqlError(result.Error) {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return execution, nil
}

func (provider ExecutionProvider) Create(
	id uuid.UUID,
	jobId uuid.UUID,
	statusType string,
	finishedAt *time.Time,
) (*model.Execution, error) {
	execution := model.Execution{
		Id:         id,
		JobId:      jobId,
		Status:     statusType,
		FinishedAt: finishedAt,
	}

	result := provider.Database.Create(&execution)
	if isSqlError(result.Error) {
		return nil, result.Error
	}
	return &execution, nil
}

func (provider ExecutionProvider) Update(id uuid.UUID, jobId *uuid.UUID, statusType *string, finishedAt *time.Time) (*model.Execution, error) {
	execution, err := provider.Get(id)
	if err != nil {
		return nil, err
	}
	if execution == nil {
		return nil, nil
	}

	if finishedAt != nil {
		execution.FinishedAt = finishedAt
	}
	if jobId != nil {
		execution.JobId = *jobId
	}
	if statusType != nil {
		execution.Status = *statusType
	}

	result := provider.Database.Save(&execution)
	if isSqlError(result.Error) {
		return nil, result.Error
	}
	return execution, nil
}

func (provider ExecutionProvider) CreateExecutionHistory(executionId uuid.UUID, percentage float64, statusType string, message string) (*model.ExecutionHistory, error) {
	executionHistory := model.ExecutionHistory{
		Id:          uuid.New(),
		ExecutionId: executionId,
		Percentage:  percentage,
		StatusType:  statusType,
		Message:     message,
	}
	result := provider.Database.Create(&executionHistory)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return &executionHistory, nil
}

func (provider ExecutionProvider) ListPaged(pagination *Pagination[model.Execution], jobId *uuid.UUID) (*Pagination[model.Execution], error) {
	var context *gorm.DB
	if jobId == nil {
		context = paginate(model.Execution{}, "id", pagination, provider.Database)
	} else {
		context = paginateWithFilter(model.Execution{}, "id", pagination, provider.Database, "job_id = ?", *jobId)
	}
	var backups []model.Execution
	result := context.Find(&backups)

	if isSqlError(result.Error) {
		return nil, result.Error
	}

	pagination.Rows = backups
	return pagination, nil
}

func (provider ExecutionProvider) ListHistory(executionId uuid.UUID) ([]*model.ExecutionHistory, error) {
	var executionHistory []*model.ExecutionHistory
	result := provider.Database.Find(&executionHistory, "execution_id = ?", executionId)

	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return executionHistory, nil
}
