package migration

import (
	"gorm.io/gorm"
)

type V2Storage struct {
}

func (V2Storage) GetName() string {
	return "v2__storage"
}

func (V2Storage) Migrate(database *gorm.DB) error {
	err := database.Exec(`ALTER TABLE storages ADD description TEXT`).Error
	if err != nil {
		return err
	}

	return nil
}
