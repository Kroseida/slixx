package backup

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/internal/satellite/application"
	"kroseida.org/slixx/internal/satellite/syncdata"
	_schedule "kroseida.org/slixx/pkg/schedule"
)

var activeSchedules = map[uuid.UUID]_schedule.Kind{}

func InitializeScheduler() {
	for id, activeSchedule := range activeSchedules {
		err := activeSchedule.Deactivate()
		if err != nil {
			// Try to deactivate other schedules even if one fails
			application.Logger.Error("Error while deactivating schedule " + id.String() + ": " + err.Error())
		}
	}

	for _, schedule := range syncdata.Container.JobSchedules {
		kind := _schedule.ValueOf(schedule.Kind)
		if kind == nil {
			// Try to initialize other schedules even if one fails
			application.Logger.Error("Error while initializing schedule " + schedule.Id.String() + ": unknown schedule kind " + schedule.Kind)
			continue
		}
		parsedConfiguration, err := kind.Parse(schedule.Configuration)
		if err != nil {
			// Try to initialize other schedules even if one fails
			application.Logger.Error("Error while initializing schedule " + schedule.Id.String() + ": " + err.Error())
			continue
		}
		kind.Initialize(parsedConfiguration, func() {
			err := Execute(uuid.New(), schedule.JobId)
			if err != nil {
				// Error while executing schedule should not stop the scheduler, also not executing the schedule again
				application.Logger.Error("Error while executing schedule " + schedule.Id.String() + ": " + err.Error())
			}
		})
		activeSchedules[schedule.Id] = kind
	}
}
