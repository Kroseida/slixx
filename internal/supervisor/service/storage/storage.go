package storage

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/internal/supervisor/datasource"
	"kroseida.org/slixx/internal/supervisor/datasource/provider"
	"kroseida.org/slixx/internal/supervisor/syncnetwork/action"
	"kroseida.org/slixx/pkg/model"
)

func Get(id uuid.UUID) (*model.Storage, error) {
	return datasource.StorageProvider.GetStorage(id)
}

func GetPaged(pagination *provider.Pagination[model.Storage]) (*provider.Pagination[model.Storage], error) {
	return datasource.StorageProvider.GetStoragesPaged(pagination)
}

func Create(name string, description string, kindName string, configuration string) (*model.Storage, error) {
	storage, err := datasource.StorageProvider.CreateStorage(name, description, kindName, configuration)
	if err != nil {
		return nil, err
	}

	go action.SyncStorages()

	return storage, nil
}

func Update(id uuid.UUID, name *string, description *string, kindName *string, configuration *string) (*model.Storage, error) {
	storage, err := datasource.StorageProvider.UpdateStorage(
		id,
		name,
		description,
		kindName,
		configuration,
	)

	if err != nil {
		return nil, err
	}

	go action.SyncStorages()

	return storage, nil
}

func Delete(id uuid.UUID) (*model.Storage, error) {
	storage, err := datasource.StorageProvider.DeleteStorage(id)
	if err != nil {
		return nil, err
	}

	go action.SyncStorages()

	return storage, nil
}
