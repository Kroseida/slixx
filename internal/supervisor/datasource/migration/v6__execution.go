package migration

import (
	"gorm.io/gorm"
)

type V6Execution struct {
}

func (V6Execution) GetName() string {
	return "v6__execution"
}

func (V6Execution) Migrate(database *gorm.DB) error {
	err := database.Exec(`CREATE TABLE executions (id char(36) NOT NULL, 
													job_id char(36) NOT NULL,
													status text NOT NULL,
													kind text NOT NULL,
                                                    created_at DATETIME, 
                                                    finished_at DATETIME,
                                                    updated_at DATETIME default NULL,
                                                    deleted_at DATETIME,
                                                    PRIMARY KEY (id)
                                                   )`).Error
	if err != nil {
		return err
	}

	err = database.Exec(`CREATE TABLE execution_histories (id char(36) NOT NULL,
														   execution_id char(36) NOT NULL,
														   percentage float NOT NULL,
														   status_type text NOT NULL,
														   message text NOT NULL,
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
