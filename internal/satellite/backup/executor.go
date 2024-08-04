package backup

import (
	"errors"
	"github.com/google/uuid"
	"kroseida.org/slixx/internal/satellite/application"
	"kroseida.org/slixx/internal/satellite/syncdata"
	"kroseida.org/slixx/internal/satellite/syncnetwork/action"
	"kroseida.org/slixx/internal/satellite/syncnetwork/manager"
	"kroseida.org/slixx/pkg/statustype"
	strategyRegistry "kroseida.org/slixx/pkg/strategy"
	"kroseida.org/slixx/pkg/utils"
	"kroseida.org/slixx/pkg/utils/parallel"
	"time"
)

var runningJobs = map[uuid.UUID]*parallel.RunningJob{}

func WatchRunningJobs() {
	for {
		for id, runningJob := range runningJobs {
			if runningJob.StartedAt.Add(time.Duration(application.CurrentSettings.Backup.Timeout) * time.Hour).Before(time.Now()) {
				runningJob.Canceled = true
				application.Logger.Info("Job with id", runningJob.JobId, "timed out and was canceled automatically")
			}
			if runningJob.Canceled {
				runningJob.Callback(utils.StatusUpdate{
					StatusType: statustype.Cancelled,
					Message:    "Job was canceled",
					Percentage: 0,
				})

				delete(runningJobs, id)
			}
		}

		time.Sleep(5 * time.Second)
	}
}

func Execute(id uuid.UUID, jobId uuid.UUID) error {
	runningJobs[id] = &parallel.RunningJob{
		JobId:     jobId,
		Canceled:  false,
		StartedAt: time.Now(),
		Callback: func(status utils.StatusUpdate) {
			application.Logger.Info("Status update", status.Message, "P", status.Percentage, status.StatusType)
			status.Id = id
			status.JobId = &jobId
			action.SendStatusUpdate(id, "BACKUP", status)
		},
	}
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
	backupInfo, err := strategy.Execute(runningJobs[id], originStorage, destinationStorage)
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

func Restore(id uuid.UUID, jobId uuid.UUID, backupId uuid.UUID) error {
	runningJobs[id] = &parallel.RunningJob{
		JobId:     jobId,
		Canceled:  false,
		StartedAt: time.Now(),
		Callback: func(status utils.StatusUpdate) {
			application.Logger.Info("Status update", status.Message, "P", status.Percentage, status.StatusType)
			status.Id = id
			status.JobId = &jobId
			action.SendStatusUpdate(id, "RESTORE", status)
		},
	}

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
	err = strategy.Restore(runningJobs[id], originStorage, destinationStorage, &backupId)
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

func Delete(id uuid.UUID, jobId uuid.UUID, backupId uuid.UUID) error {
	runningJobs[id] = &parallel.RunningJob{
		JobId:     jobId,
		Canceled:  false,
		StartedAt: time.Now(),
		Callback: func(status utils.StatusUpdate) {
			application.Logger.Info("Status update", status.Message, "P", status.Percentage, status.StatusType)
			status.Id = id
			status.JobId = &jobId
			action.SendStatusUpdate(id, "DELETE", status)
		},
	}

	application.Logger.Info("Deleting backup", backupId)
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
	destinationStorage, err := loadAndInitializeStorage(*syncdata.Container.Storages[job.DestinationStorageId])
	if err != nil {
		return err
	}

	// Execute strategy
	err = strategy.Delete(runningJobs[id], destinationStorage, &backupId)
	if err != nil {
		return err
	}

	// Close everything
	destinationStorage.Close()
	strategy.Close()

	action.SendDeleteInfo(id, jobId, backupId)
	application.Logger.Info("Backup deleted", backupId)
	return nil
}
