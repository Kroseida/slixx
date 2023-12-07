package backup

import (
	"errors"
	"github.com/google/uuid"
	"kroseida.org/slixx/internal/satellite/application"
	"kroseida.org/slixx/internal/satellite/syncdata"
	"kroseida.org/slixx/internal/satellite/syncnetwork/action"
	"kroseida.org/slixx/internal/satellite/syncnetwork/manager"
	strategyRegistry "kroseida.org/slixx/pkg/strategy"
)

func Execute(id *uuid.UUID, jobId uuid.UUID) error {
	application.Logger.Info("Executing job", jobId)
	job := syncdata.Container.Jobs[jobId]
	if job == nil {
		return errors.New("job not found")
	}
	if job.ExecutorSatelliteId != *manager.Server.Id {
		return nil
	}

	strategy := strategyRegistry.ValueOf(job.Strategy)
	if strategy == nil {
		return errors.New("strategy not found")
	}
	parsedConfiguration, err := strategy.Parse(job.Configuration)
	if err != nil {
		return err
	}

	// Initialize strategy
	err = strategy.Initialize(parsedConfiguration)
	if err != nil {
		return err
	}

	// Load storages
	originStorage, err := loadAndInitializeStorage(*syncdata.Container.Storages[job.OriginStorageId])
	if err != nil {
		return err
	}
	destinationStorage, err := loadAndInitializeStorage(*syncdata.Container.Storages[job.DestinationStorageId])
	if err != nil {
		return err
	}

	// Execute strategy
	backupInfo, err := strategy.Execute(jobId, originStorage, destinationStorage, func(status strategyRegistry.StatusUpdate) {
		application.Logger.Info("Status update", status.Message, "P", status.Percentage, status.StatusType)
		status.Id = id
		status.JobId = &job.Id
		action.SendStatusUpdate(id, "BACKUP", status)
	})
	if err != nil {
		return err
	}

	action.SendRawBackupInfo(
		backupInfo.Id,
		backupInfo.JobId,
		id,
		backupInfo.CreatedAt,
		backupInfo.OriginKind,
		backupInfo.DestinationKind,
		backupInfo.Strategy,
	)

	// Close everything
	originStorage.Close()
	destinationStorage.Close()
	strategy.Close()

	application.Logger.Info("Job executed", jobId)

	return nil
}

func Restore(id *uuid.UUID, jobId uuid.UUID, backupId uuid.UUID) error {
	application.Logger.Info("Restoring backup", backupId)
	job := syncdata.Container.Jobs[jobId]
	if job == nil {
		return errors.New("job not found")
	}
	if job.ExecutorSatelliteId != *manager.Server.Id {
		return nil
	}

	strategy := strategyRegistry.ValueOf(job.Strategy)
	if strategy == nil {
		return errors.New("strategy not found")
	}
	parsedConfiguration, err := strategy.Parse(job.Configuration)
	if err != nil {
		return err
	}

	// Initialize strategy
	err = strategy.Initialize(parsedConfiguration)
	if err != nil {
		return err
	}

	// Load storages
	originStorage, err := loadAndInitializeStorage(*syncdata.Container.Storages[job.OriginStorageId])
	if err != nil {
		return err
	}
	destinationStorage, err := loadAndInitializeStorage(*syncdata.Container.Storages[job.DestinationStorageId])
	if err != nil {
		return err
	}

	// Execute strategy
	err = strategy.Restore(originStorage, destinationStorage, &backupId, func(status strategyRegistry.StatusUpdate) {
		application.Logger.Info("Status update", status.Message, "P", status.Percentage, status.StatusType)
		status.Id = id
		status.JobId = &job.Id
		action.SendStatusUpdate(id, "RESTORE", status)
	})
	if err != nil {
		return err
	}

	// Close everything
	originStorage.Close()
	destinationStorage.Close()
	strategy.Close()

	application.Logger.Info("Backup restored", backupId)
	return nil
}
