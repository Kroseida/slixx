package backup

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/internal/supervisor/application"
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

func Restore(backupId uuid.UUID) (*uuid.UUID, error) {
	backup, err := datasource.BackupProvider.Get(backupId)
	if err != nil {
		return nil, err
	}

	backupRes, err := action.SendExecuteRestore(backup.JobId, backupId)
	if err != nil {
		return nil, err
	}
	return backupRes, nil
}

func ResyncSatellite(satelliteId uuid.UUID) error {
	action.SyncStorages(&satelliteId)
	action.SyncJobs(&satelliteId)
	return action.SendRequestBackupSync(&satelliteId)
}

func ApplyBackupToIndex(
	id uuid.UUID,
	jobId uuid.UUID,
	executionId uuid.UUID,
	createdAt time.Time,
	originKind string,
	destinationKind string,
	strategy string,
) error {
	job, _ := datasource.JobProvider.Get(jobId)

	jobName := jobId.String()
	if job != nil {
		jobName = job.Name
	}

	existingBackup, err := datasource.BackupProvider.Get(id)

	if err != nil {
		return err
	}
	if existingBackup != nil && existingBackup.Id == id {
		application.Logger.Warn("backup ", id, " is already indexed! skipping...")
		return nil
	}

	_, err = datasource.BackupProvider.Create(
		id,
		jobName+" at "+createdAt.Format("2006-01-02 15:04:05"),
		"",
		jobId,
		executionId,
		createdAt,
		originKind,
		destinationKind,
		strategy,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetPaged(pagination *provider.Pagination[model.Backup], jobId *uuid.UUID) (*provider.Pagination[model.Backup], error) {
	return datasource.BackupProvider.ListPaged(pagination, jobId)
}

func Get(id uuid.UUID) (*model.Backup, error) {
	return datasource.BackupProvider.Get(id)
}

func RequestDelete(id uuid.UUID, jobId uuid.UUID, backupId uuid.UUID) error {
	job, _ := datasource.JobProvider.Get(jobId)
	if job == nil {
		return nil
	}

	err := action.SendRequestDeleteBackup(id, jobId, backupId)
	if err != nil {
		return err
	}
	return nil
}

func DeleteBackupFromIndex(id uuid.UUID) error {
	_, err := datasource.BackupProvider.Delete(id)
	return err
}
