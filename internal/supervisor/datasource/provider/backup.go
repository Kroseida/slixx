package provider

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"kroseida.org/slixx/internal/supervisor/application"
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
	executionId *uuid.UUID,
	createdAt time.Time,
	originKind string,
	destinationKind string,
	strategy string,
) (*model.Backup, error) {
	jobName := jobId.String()

	existingBackup, err := provider.Get(id)

	if err != nil {
		return nil, err
	}

	if existingBackup != nil && existingBackup.Id == id {
		application.Logger.Warn("backup ", id, " is already indexed! skipping...")
		return existingBackup, nil
	}

	backup := model.Backup{
		Id:              id,
		Name:            jobName + " at " + createdAt.Format("2006-01-02 15:04:05"),
		Description:     "",
		ExecutionId:     executionId,
		JobId:           jobId,
		OriginKind:      originKind,
		DestinationKind: destinationKind,
		Strategy:        strategy,
		CreatedAt:       createdAt,
		UpdatedAt:       createdAt,
	}

	result := provider.Database.Create(&backup)
	if isSqlError(result.Error) {
		return nil, result.Error
	}
	return &backup, nil
}

func (provider BackupProvider) Get(id uuid.UUID) (*model.Backup, error) {
	var backup *model.Backup
	result := provider.Database.First(&backup, id)

	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return backup, nil
}

func (provider BackupProvider) List() ([]*model.Backup, error) {
	var backups []*model.Backup
	result := provider.Database.Find(&backups)

	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return backups, nil
}

func (provider BackupProvider) ListPaged(pagination *Pagination[model.Backup], jobId *uuid.UUID) (*Pagination[model.Backup], error) {
	var context *gorm.DB
	if jobId == nil {
		context = paginate(model.Backup{}, "name", pagination, provider.Database)
	} else {
		context = paginateWithFilter(model.Backup{}, "name", pagination, provider.Database, "job_id = ?", *jobId)
	}
	var backups []model.Backup
	result := context.Find(&backups)

	if isSqlError(result.Error) {
		return nil, result.Error
	}

	pagination.Rows = backups
	return pagination, nil
}
