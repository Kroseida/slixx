package datasource

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"kroseida.org/slixx/internal/supervisor/application"
	"kroseida.org/slixx/internal/supervisor/datasource/migration"
	"kroseida.org/slixx/internal/supervisor/datasource/provider"
)

var StorageProvider provider.StorageProvider
var UserProvider provider.UserProvider
var SatelliteProvider provider.SatelliteProvider
var JobProvider provider.JobProvider
var localDatabase *gorm.DB

func Connect() error {
	var err error
	if application.CurrentSettings.Database.Kind == "sqlite" {
		err = ConnectSqlite()
	}
	if err != nil {
		return err
	}

	err = migration.Migrate(localDatabase)
	if err != nil {
		return err
	}

	StorageProvider = provider.StorageProvider{
		Database:    localDatabase,
		JobProvider: &JobProvider,
	}
	UserProvider = provider.UserProvider{Database: localDatabase}
	SatelliteProvider = provider.SatelliteProvider{
		Database: localDatabase,
	}
	JobProvider = provider.JobProvider{
		Database:          localDatabase,
		StorageProvider:   &StorageProvider,
		SatelliteProvider: &SatelliteProvider,
	}

	err = UserProvider.Init()
	if err != nil {
		return err
	}

	return nil
}

func ConnectSqlite() error {
	var logMode logger.LogLevel
	if application.CurrentSettings.Logger.Mode == "debug" {
		logMode = logger.Info
	} else {
		logMode = logger.Silent
	}

	database, err := gorm.Open(sqlite.Open(application.CurrentSettings.Database.Configuration["file"]), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logMode),
	})
	if err != nil {
		return err
	}
	localDatabase = database
	return nil
}

func Close() error {
	sqlDb, err := localDatabase.DB()
	if err != nil {
		return err
	}
	err = sqlDb.Close()
	if err != nil {
		return err
	}
	return nil
}
