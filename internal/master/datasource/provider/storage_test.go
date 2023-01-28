package provider_test

import (
	"github.com/stretchr/testify/assert"
	"kroseida.org/slixx/internal/master/datasource"
	"testing"
)

func Test_CreateStorage(t *testing.T) {
	teardownSuite := setupSuite()

	_, err := datasource.StorageProvider.CreateStorage("Test Storage", "test_kind", "{}")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	storages, err := datasource.StorageProvider.GetStorages()
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, 1, len(storages))
	assert.Equal(t, "Test Storage", storages[0].Name)
	assert.Equal(t, "test_kind", storages[0].Kind)
	assert.Equal(t, "{}", storages[0].Configuration)
	teardownSuite()
}

func Test_GetStorages(t *testing.T) {
	teardownSuite := setupSuite()

	_, err := datasource.StorageProvider.CreateStorage("Test Storage", "test_kind", "{}")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	_, err = datasource.StorageProvider.CreateStorage("Test Storage 2", "test_kind 2", "{}")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	storages, err := datasource.StorageProvider.GetStorages()
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, 2, len(storages))
	assert.Equal(t, "Test Storage", storages[0].Name)
	assert.Equal(t, "test_kind", storages[0].Kind)
	assert.Equal(t, "{}", storages[0].Configuration)
	assert.Equal(t, "Test Storage 2", storages[1].Name)
	assert.Equal(t, "test_kind 2", storages[1].Kind)
	assert.Equal(t, "{}", storages[1].Configuration)
	teardownSuite()
}

func Test_GetStorage(t *testing.T) {
	teardownSuite := setupSuite()

	_, err := datasource.StorageProvider.CreateStorage("Test Storage", "test_kind", "{}")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	createdStorage, err := datasource.StorageProvider.CreateStorage("Test Storage", "test_kind", "{}")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	storage, err := datasource.StorageProvider.GetStorage(createdStorage.Id)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, "Test Storage", storage.Name)
	assert.Equal(t, "test_kind", storage.Kind)
	assert.Equal(t, "{}", storage.Configuration)
	teardownSuite()
}
