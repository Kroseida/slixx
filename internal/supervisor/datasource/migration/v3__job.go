package migration

import (
	"gorm.io/gorm"
)

type V3Job struct {
}

func (V3Job) GetName() string {
	return "v3__job"
}

func (V3Job) Migrate(database *gorm.DB) error {
	err := database.Exec(`CREATE TABLE jobs (id char(36) NOT NULL, 
                                                    name text NOT NULL, 
                                                    description text NOT NULL,
                                                    strategy text NOT NULL,
                                                    configuration text NOT NULL,
                                                    created_at DATETIME, 
                                                    updated_at DATETIME, 
                                                    deleted_at DATETIME,
                                                    origin_storage_id char(36) NOT NULL,
                                                    destination_storage_id char(36) NOT NULL, 
                                                    PRIMARY KEY (id)
                                                   )`).Error
	if err != nil {
		return err
	}

	return nil
}
