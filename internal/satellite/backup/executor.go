package backup

import (
	"errors"
	"github.com/google/uuid"
	"kroseida.org/slixx/internal/satellite/application"
	"kroseida.org/slixx/internal/satellite/syncdata"
	"kroseida.org/slixx/internal/satellite/syncnetwork/action"
	"kroseida.org/slixx/internal/satellite/syncnetwork/manager"
	"kroseida.org/slixx/pkg/model"
	"kroseida.org/slixx/pkg/storage"
	storageRegistry "kroseida.org/slixx/pkg/storage"
	strategyRegistry "kroseida.org/slixx/pkg/strategy"
)

func JobCheckupLoop() {
	for {
		for _, job := range syncdata.Container.Jobs {
			if &job.Id != manager.Server.Id {
				continue
			}
			id := uuid.New()

			go func() {
				err := Execute(&id, job.Id)
				if err != nil {
					application.Logger.Error(err)
				}
			}()
		}
	}
}

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
	err = strategy.Execute(originStorage, destinationStorage, func(status strategyRegistry.BackupStatusUpdate) {
		application.Logger.Info("Status update", status.Message, "P", status.Percentage, status.StatusType)
		status.JobId = &job.Id
		action.SendBackupStatusUpdate(id, status)
	})
	if err != nil {
		return err
	}

	// Close everything
	originStorage.Close()
	destinationStorage.Close()
	strategy.Close()

	application.Logger.Info("Job executed", jobId)

	return nil
}

func loadAndInitializeStorage(storage model.Storage) (storage.Kind, error) {
	kind := storageRegistry.ValueOf(storage.Kind)
	if kind == nil {
		return nil, errors.New("storage not found")
	}
	parsedConfiguration, err := kind.Parse(storage.Configuration)
	if err != nil {
		return nil, err
	}

	err = kind.Initialize(parsedConfiguration)
	if err != nil {
		return nil, err
	}

	return kind, nil
}
