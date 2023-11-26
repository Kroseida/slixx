package controller_test

import (
	"github.com/stretchr/testify/assert"
	"kroseida.org/slixx/internal/supervisor/datasource"
	"kroseida.org/slixx/internal/supervisor/graphql/controller"
	"testing"
)

func Test_GetStorage(t *testing.T) {
	teardownSuite := setupSuite()
	created, err := controller.CreateStorage(withPermissions([]string{"storage.create"}), controller.CreateStorageDto{
		Name:          "Testaaaaaa",
		Description:   "description",
		Kind:          "FTP",
		Configuration: "{}", // we expect to make this to a valid json
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	storage, err := controller.GetStorage(withPermissions([]string{"storage.view"}), controller.GetStorageDto{
		Id: created.Id,
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, created.Id, storage.Id)
	assert.Equal(t, "Testaaaaaa", storage.Name)
	assert.Equal(t, "description", storage.Description)
	assert.Equal(t, "FTP", storage.Kind)
	assert.Equal(t, "{\"host\":\"\",\"timeout\":0,\"file\":\"/\",\"username\":\"\",\"password\":\"\"}", storage.Configuration)

	teardownSuite()
}

func Test_GetStorages(t *testing.T) {
	teardownSuite := setupSuite()
	_, err := controller.CreateStorage(withPermissions([]string{"storage.create"}), controller.CreateStorageDto{
		Name:          "Testaaaaaa",
		Description:   "description",
		Kind:          "FTP",
		Configuration: "{}", // we expect to make this to a valid json
	})
	_, err = controller.CreateStorage(withPermissions([]string{"storage.create"}), controller.CreateStorageDto{
		Name:          "Testaaaaaa2",
		Description:   "description",
		Kind:          "FTP",
		Configuration: "{}", // we expect to make this to a valid json
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	storages, err := controller.GetStorages(withPermissions([]string{"storage.view"}), controller.PageArgs{})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, 3, len(storages.Rows))
	teardownSuite()
}

func Test_UpdateStorage(t *testing.T) {
	teardownSuite := setupSuite()
	updatedName := "Testaaaaaa2"
	created, err := controller.CreateStorage(withPermissions([]string{"storage.create"}), controller.CreateStorageDto{
		Name:          "Testaaaaaa",
		Description:   "description",
		Kind:          "FTP",
		Configuration: "{}", // we expect to make this to a valid json
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	_, err = controller.UpdateStorage(withPermissions([]string{"storage.update"}), controller.UpdateStorageDto{
		Id:   created.Id,
		Name: &updatedName,
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	storage, err := controller.GetStorage(withPermissions([]string{"storage.view"}), controller.GetStorageDto{
		Id: created.Id,
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, created.Id, storage.Id)
	assert.Equal(t, updatedName, storage.Name)

	teardownSuite()
}

func Test_DeleteStorage(t *testing.T) {
	teardownSuite := setupSuite()
	created, err := controller.CreateStorage(withPermissions([]string{"storage.create"}), controller.CreateStorageDto{
		Name:          "Testaaaaaa",
		Description:   "description",
		Kind:          "FTP",
		Configuration: "{}", // we expect to make this to a valid json
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	_, err = controller.DeleteStorage(withPermissions([]string{"storage.delete"}), controller.DeleteStorageDto{
		Id: created.Id,
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	storageInDatabase, err := datasource.UserProvider.Get(created.Id)

	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Nil(t, storageInDatabase)

	teardownSuite()
}

func Test_GetStorageKinds(t *testing.T) {
	teardownSuite := setupSuite()
	kinds, err := controller.GetStorageKinds()
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, 1, len(kinds))
	assert.Equal(t, "FTP", kinds[0].Name)
	teardownSuite()
}
