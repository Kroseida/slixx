package migration_test

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"kroseida.org/slixx/internal/supervisior/application"
	"os"
	"time"
)

func setupSuite() (func(), *gorm.DB) {
	os.Remove("slixx.db")

	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.Encoding = "console"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	logger, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}

	application.Logger = logger.Sugar()

	database, err := gorm.Open(sqlite.Open("slixx.db"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		panic(err)
	}

	return func() {
		sqlDb, closeError := database.DB()
		if closeError != nil {
			panic(closeError)
		}
		err = sqlDb.Close()
		if err != nil {
			panic(closeError)
		}
		os.Remove("slixx.db")
	}, database
}
