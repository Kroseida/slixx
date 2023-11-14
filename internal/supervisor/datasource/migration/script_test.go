package migration_test

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"kroseida.org/slixx/internal/supervisor/datasource/migration"
	"testing"
)

type TestMigration struct {
}

type Blub struct {
	gorm.Model
	Id   int
	Name string
}

func (TestMigration) GetName() string {
	return "v1__test"
}

func (TestMigration) Migrate(database *gorm.DB) error {
	err := database.Exec(`CREATE TABLE blubs (id INTEGER PRIMARY KEY AUTOINCREMENT, 
												  name TEXT, 
												  created_at DATETIME, 
												  updated_at DATETIME, 
												  deleted_at DATETIME
                   								 )`).Error
	if err != nil {
		return err
	}
	return nil
}

func Test_MigrateScripts(t *testing.T) {
	teardownSuite, database := setupSuite()
	err := migration.MigrateScripts(database, []migration.Script{
		TestMigration{},
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	err = database.Create(&Blub{
		Name: "Test",
	}).Error
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}
	var blub Blub

	result := database.First(&blub, "name = ?", "Test")
	if result.Error != nil {
		t.Error(result.Error)
		return
	}

	assert.Equal(t, int64(1), result.RowsAffected)
	assert.Equal(t, "Test", blub.Name)
	assert.Equal(t, 1, blub.Id)
	teardownSuite()
}

func Test_Migrate(t *testing.T) {
	teardownSuite, database := setupSuite()
	err := migration.Migrate(database)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}
	teardownSuite()
}
