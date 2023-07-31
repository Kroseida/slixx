package provider

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"kroseida.org/slixx/pkg/model"
	"time"
)

// BackupProvider Backup Provider
type BackupProvider struct {
	Database    *gorm.DB
	JobProvider *JobProvider
}

func (provider BackupProvider) ApplyBackupToIndex(
	id uuid.UUID,
	jobId uuid.UUID,
	executionId uuid.UUID,
	createdAt time.Time,
) (*model.Backup, error) {
	jobName := jobId.String()

	job, _ := provider.JobProvider.Get(jobId)
	if job == nil {
		jobName = job.Name
	}
	backup := model.Backup{
		Id:          id,
		Name:        jobName + " at " + createdAt.Format("2006-01-02 15:04:05"),
		Description: "",
		ExecutionId: executionId,
		JobId:       jobId,
		CreatedAt:   createdAt,
	}

	result := provider.Database.Create(backup)
	if isSqlError(result.Error) {
		return nil, result.Error
	}
	return &backup, nil
}

func (provider BackupProvider) List() ([]*model.Backup, error) {
	var backups []*model.Backup
	result := provider.Database.Find(&backups)

	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return backups, nil
}

func (provider BackupProvider) ListPaged(pagination *Pagination[model.Backup]) (*Pagination[model.Backup], error) {
	context := paginate(model.Backup{}, "name", pagination, provider.Database)

	var backups []model.Backup
	result := context.Find(&backups)

	if isSqlError(result.Error) {
		return nil, result.Error
	}

	pagination.Rows = backups
	return pagination, nil
}
