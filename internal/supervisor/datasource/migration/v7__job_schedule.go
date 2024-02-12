package migration

import (
	"gorm.io/gorm"
)

type V7JobSchedule struct {
}

func (V7JobSchedule) GetName() string {
	return "v7__job_schedule"
}

func (V7JobSchedule) Migrate(database *gorm.DB) error {
	err := database.Exec(`CREATE TABLE job_schedules (
    		id UUID PRIMARY KEY,
    		job_id UUID NOT NULL,
    		name text NOT NULL, 
        	description text NOT NULL,
    		kind VARCHAR(255) NOT NULL,
    		configuration TEXT NOT NULL,
			created_at DATETIME, 
			updated_at DATETIME default NULL,
			deleted_at DATETIME
		)`).Error
	if err != nil {
		return err
	}

	return nil
}
