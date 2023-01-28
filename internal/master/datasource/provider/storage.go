package provider

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"kroseida.org/slixx/internal/master/datasource/model"
)

// StorageProvider Storage Provider
type StorageProvider struct {
	Database *gorm.DB
}

func (provider StorageProvider) DeleteStorage(id uuid.UUID) (bool, error) {
	var storage *model.Storage
	result := provider.Database.Delete(&storage, "id = ?", id)

	if isSqlError(result.Error) {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (provider StorageProvider) CreateStorage(name string, kind string, configuration string) (*model.Storage, error) {
	storage := &model.Storage{
		Id:            uuid.New(),
		Name:          name,
		Kind:          kind,
		Configuration: configuration,
	}

	result := provider.Database.Create(&storage)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return storage, nil
}

func (provider StorageProvider) GetStorages() ([]*model.Storage, error) {
	var storages []*model.Storage
	result := provider.Database.Find(&storages)

	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return storages, nil
}

func (provider StorageProvider) GetStorage(id uuid.UUID) (*model.Storage, error) {
	var storage *model.Storage
	result := provider.Database.First(&storage, "id = ?", id)

	if isSqlError(result.Error) {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return storage, nil
}
