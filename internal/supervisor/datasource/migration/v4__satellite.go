package migration

import (
	"gorm.io/gorm"
)

type V4Satellite struct {
}

func (V4Satellite) GetName() string {
	return "v4__satellite"
}

func (V4Satellite) Migrate(database *gorm.DB) error {
	err := database.Exec(`CREATE TABLE satellites (id char(36) NOT NULL, 
                                                    name text NOT NULL, 
                                                    description text NOT NULL,
													address text NOT NULL,
													token text NOT NULL,
                                                    created_at DATETIME, 
                                                    updated_at DATETIME, 
                                                    deleted_at DATETIME,
                                                    PRIMARY KEY (id)
                                                   )`).Error
	if err != nil {
		return err
	}

	err = database.Exec(`CREATE TABLE satellite_log_entries (id char(36) NOT NULL,
														  satellite_id char(36) NOT NULL,
														  sender text NOT NULL,
														  message text NOT NULL,
														  level char(36) NOT NULL,
														  logged_at DATETIME,									  
														  created_at DATETIME,
														  updated_at DATETIME,
														  deleted_at DATETIME,
														  PRIMARY KEY (id)
														 )`).Error
	if err != nil {
		return err
	}

	err = database.Exec(`ALTER TABLE jobs ADD COLUMN executor_satellite_id char(36) DEFAULT NULL`).Error
	if err != nil {
		return err
	}

	return nil
}
