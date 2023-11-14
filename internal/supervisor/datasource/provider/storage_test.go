package provider_test

import (
	"github.com/stretchr/testify/assert"
	"kroseida.org/slixx/internal/supervisor/datasource"
	"kroseida.org/slixx/internal/supervisor/datasource/provider"
	"kroseida.org/slixx/pkg/model"
	"testing"
)

func Test_DeleteStorage(t *testing.T) {
	teardownSuite := setupSuite()

	createdStorage, err := datasource.StorageProvider.Create("Test Storage", "", "FTP", "{}")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	_, err = datasource.StorageProvider.Delete(createdStorage.Id)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	storages, err := datasource.StorageProvider.List()
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, 0, len(storages))
	teardownSuite()
}

func Test_CreateStorage_EmptyName(t *testing.T) {
	teardownSuite := setupSuite()

	_, err := datasource.StorageProvider.Create("", "", "FTP", "{}")
	if err == nil {
		t.Error("Expected error")
		teardownSuite()
		return
	}

	assert.Equal(t, "name can not be empty", err.Error())
	teardownSuite()
}

func Test_CreateStorage_InvalidKind(t *testing.T) {
	teardownSuite := setupSuite()

	_, err := datasource.StorageProvider.Create("Test Storage", "", "INVALID", "{}")
	if err == nil {
		t.Error("Expected error")
		teardownSuite()
		return
	}

	assert.Equal(t, "Invalid storage kind \"INVALID\"", err.Error())
	teardownSuite()
}

func Test_CreateStorage(t *testing.T) {
	teardownSuite := setupSuite()

	_, err := datasource.StorageProvider.Create("Test Storage", "", "FTP", "{}")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	storages, err := datasource.StorageProvider.List()
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, 1, len(storages))
	assert.Equal(t, "Test Storage", storages[0].Name)
	assert.Equal(t, "FTP", storages[0].Kind)
	assert.Equal(t, "{\"host\":\"\",\"timeout\":0,\"file\":\"\",\"username\":\"\",\"password\":\"\"}", storages[0].Configuration)
	teardownSuite()
}

func Test_UpdateStorage_InvalidKind(t *testing.T) {
	teardownSuite := setupSuite()

	newKind := "Invalid"
	createdStorage, err := datasource.StorageProvider.Create("Test Storage", "", "FTP", "{}")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	_, err = datasource.StorageProvider.Update(createdStorage.Id, nil, nil, &newKind, nil)
	if err == nil {
		t.Error("Expected error")
		teardownSuite()
		return
	}
	teardownSuite()
}

func Test_UpdateStorage(t *testing.T) {
	teardownSuite := setupSuite()

	newName := "Updated Storage"
	newConfiguration := "{}"
	createdStorage, err := datasource.StorageProvider.Create("Test Storage", "", "FTP", "{}")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	_, err = datasource.StorageProvider.Update(createdStorage.Id, &newName, nil, nil, &newConfiguration)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	updatedStorage, err := datasource.StorageProvider.Get(createdStorage.Id)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, "Updated Storage", updatedStorage.Name)
	assert.Equal(t, "{\"host\":\"\",\"timeout\":0,\"file\":\"\",\"username\":\"\",\"password\":\"\"}", updatedStorage.Configuration)
	teardownSuite()
}

func Test_UpdateStorage_EmptyName(t *testing.T) {
	teardownSuite := setupSuite()

	newName := ""
	createdStorage, err := datasource.StorageProvider.Create("Test Storage", "", "FTP", "{}")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	_, err = datasource.StorageProvider.Update(createdStorage.Id, &newName, nil, nil, nil)
	assert.NotNil(t, err)
	teardownSuite()
}

func Test_GetStorages(t *testing.T) {
	teardownSuite := setupSuite()

	_, err := datasource.StorageProvider.Create("Test Storage", "", "FTP", "{}")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	_, err = datasource.StorageProvider.Create("Test Storage 2", "", "FTP", "{}")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	storages, err := datasource.StorageProvider.List()
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, 2, len(storages))
	assert.Equal(t, "Test Storage", storages[0].Name)
	assert.Equal(t, "FTP", storages[0].Kind)
	assert.Equal(t, "{\"host\":\"\",\"timeout\":0,\"file\":\"/\",\"username\":\"\",\"password\":\"\"}", storages[0].Configuration)
	assert.Equal(t, "Test Storage 2", storages[1].Name)
	assert.Equal(t, "FTP", storages[1].Kind)
	assert.Equal(t, "{\"host\":\"\",\"timeout\":0,\"file\":\"/\",\"username\":\"\",\"password\":\"\"}", storages[1].Configuration)
	teardownSuite()
}

func Test_GetStoragesPaged(t *testing.T) {
	teardownSuite := setupSuite()

	_, err := datasource.StorageProvider.Create("Test Storage", "", "FTP", "{}")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	_, err = datasource.StorageProvider.Create("Test Storage 2", "", "FTP", "{}")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	storages, err := datasource.StorageProvider.ListPaged(&provider.Pagination[model.Storage]{
		Page:  1,
		Limit: 1,
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, 1, len(storages.Rows))
	teardownSuite()
}

func Test_GetStorage(t *testing.T) {
	teardownSuite := setupSuite()

	_, err := datasource.StorageProvider.Create("Test Storage", "", "FTP", "{}")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	createdStorage, err := datasource.StorageProvider.Create("Test Storage", "", "FTP", "{}")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	storage, err := datasource.StorageProvider.Get(createdStorage.Id)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, "Test Storage", storage.Name)
	assert.Equal(t, "FTP", storage.Kind)
	assert.Equal(t, "{\"host\":\"\",\"timeout\":0,\"file\":\"/\",\"username\":\"\",\"password\":\"\"}", storage.Configuration)
	teardownSuite()
}
