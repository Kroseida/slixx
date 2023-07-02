package storageservice

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/internal/supervisor/datasource"
	"kroseida.org/slixx/internal/supervisor/syncnetwork/action"
	"kroseida.org/slixx/pkg/model"
)

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
