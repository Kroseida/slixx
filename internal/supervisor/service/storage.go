package service

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/internal/supervisor/datasource"
	"kroseida.org/slixx/internal/supervisor/syncnetwork"
	"kroseida.org/slixx/pkg/model"
)

func CreateStorage(name string, description string, kindName string, configuration string) (*model.Storage, error) {
	storage, err := datasource.StorageProvider.CreateStorage(name, description, kindName, configuration)
	if err != nil {
		return nil, err
	}

	go syncnetwork.SyncStorages()

	return storage, nil
}

func UpdateStorage(id uuid.UUID, name *string, description *string, kindName *string, configuration *string) (*model.Storage, error) {
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

	go syncnetwork.SyncStorages()

	return storage, nil
}

func DeleteStorage(id uuid.UUID) (*model.Storage, error) {
	storage, err := datasource.StorageProvider.DeleteStorage(id)
	if err != nil {
		return nil, err
	}

	go syncnetwork.SyncStorages()

	return storage, nil
}
