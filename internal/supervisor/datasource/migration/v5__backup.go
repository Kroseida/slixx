package migration

import (
	"gorm.io/gorm"
)

type V5Backup struct {
}

func (V5Backup) GetName() string {
	return "v5__backup"
}

func (V5Backup) Migrate(database *gorm.DB) error {
	err := database.Exec(`CREATE TABLE backups (id char(36) NOT NULL, 
                                                    name text NOT NULL, 
                                                    description text NOT NULL,
													job_id char(36) NOT NULL,
													execution_id char(36),
                                                    created_at DATETIME, 
                                                    updated_at DATETIME, 
                                                    deleted_at DATETIME,
                                                    PRIMARY KEY (id)
                                                   )`).Error
	if err != nil {
		return err
	}

	return nil
}
