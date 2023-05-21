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

	err = database.Exec(`CREATE TABLE satellite_logs (id char(36) NOT NULL,
														  satellite_id char(36) NOT NULL,
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
