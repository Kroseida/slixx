package backup

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/internal/supervisor/datasource"
	"kroseida.org/slixx/internal/supervisor/datasource/provider"
	"kroseida.org/slixx/internal/supervisor/syncnetwork/action"
	"kroseida.org/slixx/pkg/model"
	"time"
)

func Execute(jobId uuid.UUID) {
	action.SendExecuteBackup(jobId)
}

func ApplyBackupToIndex(
	id uuid.UUID,
	jobId uuid.UUID,
	executionId uuid.UUID,
	createdAt time.Time,
) (*model.Backup, error) {
	backup, err := datasource.BackupProvider.ApplyBackupToIndex(
		id,
		jobId,
		executionId,
		createdAt,
	)
	if err != nil {
		return nil, err
	}
	return backup, nil
}

func GetPaged(pagination *provider.Pagination[model.Backup]) (*provider.Pagination[model.Backup], error) {
	return datasource.BackupProvider.ListPaged(pagination)
}
