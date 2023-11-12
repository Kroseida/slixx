package backup

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/internal/supervisor/datasource"
	"kroseida.org/slixx/internal/supervisor/datasource/provider"
	"kroseida.org/slixx/internal/supervisor/syncnetwork/action"
	"kroseida.org/slixx/pkg/model"
	"time"
)

func Execute(jobId uuid.UUID) (*uuid.UUID, error) {
	backup, err := action.SendExecuteBackup(jobId)
	if err != nil {
		return nil, err
	}
	return backup, nil
}

func ResyncSatellite(satelliteId uuid.UUID) error {
	action.SyncStorages(&satelliteId)
	action.SyncJobs(&satelliteId)
	return action.SendRequestBackupSync(satelliteId)
}

func ApplyBackupToIndex(
	id uuid.UUID,
	jobId uuid.UUID,
	executionId *uuid.UUID,
	createdAt time.Time,
	originKind string,
	destinationKind string,
	strategy string,
) (*model.Backup, error) {
	backup, err := datasource.BackupProvider.ApplyBackupToIndex(
		id,
		jobId,
		executionId,
		createdAt,
		originKind,
		destinationKind,
		strategy,
	)
	if err != nil {
		return nil, err
	}
	return backup, nil
}

func GetPaged(pagination *provider.Pagination[model.Backup], jobId *uuid.UUID) (*provider.Pagination[model.Backup], error) {
	return datasource.BackupProvider.ListPaged(pagination, jobId)
}
