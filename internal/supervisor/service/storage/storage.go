package storage

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/samsarahq/thunder/graphql"
	"kroseida.org/slixx/internal/supervisor/datasource"
	"kroseida.org/slixx/internal/supervisor/datasource/provider"
	"kroseida.org/slixx/internal/supervisor/syncnetwork/action"
	"kroseida.org/slixx/pkg/model"
	_storage "kroseida.org/slixx/pkg/storage"
	"reflect"
)

func Get(id uuid.UUID) (*model.Storage, error) {
	storage, err := datasource.StorageProvider.Get(id)
	if err != nil {
		return nil, err
	}
	hideSensitiveData(storage)
	return storage, nil
}

func GetPaged(pagination *provider.Pagination[model.Storage]) (*provider.Pagination[model.Storage], error) {
	page, err := datasource.StorageProvider.ListPaged(pagination)
	if err != nil {
		return nil, err
	}
	for i := range page.Rows {
		hideSensitiveData(&page.Rows[i])
	}
	return page, nil
}

func Create(name string, description string, kindName string, configuration string) (*model.Storage, error) {
	storage, err := datasource.StorageProvider.Create(name, description, kindName, configuration)
	if err != nil {
		return nil, err
	}

	go action.SyncStorages(nil)
	hideSensitiveData(storage)

	return storage, nil
}

func Update(id uuid.UUID, name *string, description *string, kindName *string, configuration *string) (*model.Storage, error) {
	storage, err := datasource.StorageProvider.Update(
		id,
		name,
		description,
		kindName,
		configuration,
	)

	if err != nil {
		return nil, err
	}

	go action.SyncStorages(nil)
	hideSensitiveData(storage)

	return storage, nil
}

func Delete(id uuid.UUID) (*model.Storage, error) {
	jobs, err := datasource.JobProvider.GetByStorageId(id)

	if err != nil {
		return nil, err
	}
	if len(jobs) > 0 {
		return nil, graphql.NewSafeError("storage is in use")
	}

	storage, err := datasource.StorageProvider.Delete(id)
	if err != nil {
		return nil, err
	}

	go action.SyncStorages(nil)
	hideSensitiveData(storage)

	return storage, nil
}

func hideSensitiveData(storage *model.Storage) error {
	configuration, err := _storage.ValueOf(storage.Kind).Parse(storage.Configuration)

	if err != nil {
		return err
	}

	val := reflect.ValueOf(configuration).Elem()
	for i := 0; i < val.NumField(); i++ {
		if val.Type().Field(i).Tag.Get("slixx") == "TOKEN" || val.Type().Field(i).Tag.Get("slixx") == "PASSWORD" {
			val.Field(i).SetString("********")
		}
	}
	bytesConfiguration, err := json.Marshal(configuration)
	if err != nil {
		return err
	}
	storage.Configuration = string(bytesConfiguration)
	return nil
}
