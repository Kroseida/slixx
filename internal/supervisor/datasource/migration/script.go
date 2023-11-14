package migration

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"kroseida.org/slixx/internal/supervisor/application"
	"kroseida.org/slixx/pkg/model"
)

type Script interface {
	GetName() string
	Migrate(db *gorm.DB) error
}

var migrationScripts = []Script{
	V1Initialize{},
	V2Storage{},
	V3Job{},
	V4Satellite{},
	V5Backup{},
}

func Migrate(database *gorm.DB) error {
	return MigrateScripts(database, migrationScripts)
}

func MigrateScripts(database *gorm.DB, scripts []Script) error {
	if (!database.Migrator().HasTable(&model.Migration{})) {
		database.Migrator().CreateTable(&model.Migration{})
	}
	for migration := range scripts {
		err := database.Transaction(func(context *gorm.DB) error {
			targetMigration := scripts[migration]

			// Check if migration is already applied
			var existingMigration model.Migration
			context.First(&existingMigration, "name = ?", targetMigration.GetName())
			if existingMigration != (model.Migration{}) {
				return nil
			}

			application.Logger.Info("Migrating: " + scripts[migration].GetName())
			context.Create(&model.Migration{
				Id:   uuid.New(),
				Name: targetMigration.GetName(),
			})
			return targetMigration.Migrate(context)
		})
		if err != nil {
			return err
		}
	}
	return nil
}
