package backup

import (
	"errors"
	"github.com/google/uuid"
	"kroseida.org/slixx/internal/satellite/application"
	"kroseida.org/slixx/internal/satellite/syncdata"
	"kroseida.org/slixx/internal/satellite/syncnetwork/action"
	"kroseida.org/slixx/pkg/model"
	"kroseida.org/slixx/pkg/storage"
	storageRegistry "kroseida.org/slixx/pkg/storage"
	strategyRegistry "kroseida.org/slixx/pkg/strategy"
)

func SendBackupInfos() {
	for _, job := range syncdata.Container.Jobs {
		strategy := strategyRegistry.ValueOf(job.Strategy)
		if strategy == nil {
			application.Logger.Error("Unknown strategy of job", job.Id, "(", job.Strategy, ")")
			continue
		}
		parsedConfiguration, err := strategy.Parse(job.Configuration)
		if err != nil {
			application.Logger.Error("Error while listing backups of job("+job.Id.String()+")", err)
			continue
		}

		// Initialize strategy
		err = strategy.Initialize(parsedConfiguration)
		if err != nil {
			application.Logger.Error("Error while listing backups of job("+job.Id.String()+")", err)
			continue
		}

		destinationStorage, err := loadAndInitializeStorage(*syncdata.Container.Storages[job.DestinationStorageId])
		if err != nil {
			application.Logger.Error("Error while listing backups of job("+job.Id.String()+")", err)
			continue
		}

		// Get backup infos
		backupInfos, err := strategy.ListBackups(destinationStorage)
		if err != nil {
			application.Logger.Error("Error while listing backups of job("+job.Id.String()+")", err)
			continue
		}

		// Send backup infos so supervisor can update its database
		for _, backupInfo := range backupInfos {
			action.SendRawBackupInfo(
				backupInfo.Id,
				&job.Id,
				uuid.UUID{},
				backupInfo.CreatedAt,
				backupInfo.OriginKind,
				backupInfo.DestinationKind,
				backupInfo.Strategy,
			)
		}

		// Close everything
		err = destinationStorage.Close()
		if err != nil {
			application.Logger.Error("Error while listing backups of job("+job.Id.String()+")", err)
			continue
		}
		err = strategy.Close()
		if err != nil {
			application.Logger.Error("Error while listing backups of job("+job.Id.String()+")", err)
			continue
		}
	}
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