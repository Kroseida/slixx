package provider

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/samsarahq/thunder/graphql"
	"gorm.io/gorm"
	"kroseida.org/slixx/pkg/model"
	_storage "kroseida.org/slixx/pkg/storage"
)

// StorageProvider Storage Provider
type StorageProvider struct {
	Database *gorm.DB
}

func (provider StorageProvider) Delete(id uuid.UUID) (*model.Storage, error) {
	storage, err := provider.Get(id)
	if storage == nil {
		return nil, graphql.NewSafeError("storage not found")
	}
	if err != nil {
		return nil, err
	}

	result := provider.Database.Delete(&storage)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return storage, nil
}

func (provider StorageProvider) Create(name string, description string, kindName string, configuration string) (*model.Storage, error) {
	if name == "" {
		return nil, graphql.NewSafeError("name can not be empty")
	}
	kind := _storage.ValueOf(kindName)
	if kind == nil {
		return nil, graphql.NewSafeError("Invalid storage kind \"%s\"", kindName)
	}
	// Check if configuration is valid
	parsedConfiguration, err := kind.Parse(configuration)
	if err != nil {
		return nil, err
	}

	rawConfiguration, err := json.Marshal(parsedConfiguration)
	if err != nil {
		return nil, err
	}

	configuration = string(rawConfiguration)

	storage := model.Storage{
		Id:            uuid.New(),
		Name:          name,
		Description:   description,
		Kind:          kindName,
		Configuration: string(rawConfiguration),
	}

	result := provider.Database.Create(&storage)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return &storage, nil
}

func (provider StorageProvider) Update(id uuid.UUID, name *string, description *string, kindName *string, configuration *string) (*model.Storage, error) {
	updateStorage, err := provider.Get(id)
	if updateStorage == nil {
		return nil, graphql.NewSafeError("storage not found")
	}
	if err != nil {
		return nil, err
	}

	if name != nil {
		if *name == "" {
			return nil, graphql.NewSafeError("name can not be empty")
		}
		updateStorage.Name = *name
	}
	if kindName != nil {
		kind := _storage.ValueOf(*kindName)
		if kind == nil {
			return nil, graphql.NewSafeError("Invalid storage kind \"%s\"", *kindName)
		}
		updateStorage.Kind = *kindName
	}
	if configuration != nil {
		kindType := _storage.ValueOf(updateStorage.Kind)

		parsedConfiguration, err := kindType.Parse(*configuration)
		if err != nil {
			return nil, err
		}

		rawConfiguration, err := json.Marshal(parsedConfiguration)
		if err != nil {
			return nil, err
		}

		updateStorage.Configuration = string(rawConfiguration)
	}
	if description != nil {
		updateStorage.Description = *description
	}

	result := provider.Database.Save(&updateStorage)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return updateStorage, nil
}

func (provider StorageProvider) List() ([]*model.Storage, error) {
	var storages []*model.Storage
	result := provider.Database.Find(&storages)

	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return storages, nil
}

func (provider StorageProvider) ListPaged(pagination *Pagination[model.Storage]) (*Pagination[model.Storage], error) {
	context := paginate(model.Storage{}, "name", pagination, provider.Database)

	var storages []model.Storage
	result := context.Find(&storages)

	if isSqlError(result.Error) {
		return nil, result.Error
	}

	pagination.Rows = storages
	return pagination, nil
}

func (provider StorageProvider) Get(id uuid.UUID) (*model.Storage, error) {
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
