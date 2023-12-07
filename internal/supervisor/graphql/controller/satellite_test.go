package controller_test

import (
	"github.com/stretchr/testify/assert"
	"kroseida.org/slixx/internal/supervisor/graphql/controller"
	"testing"
)

func Test_GetSatellite(t *testing.T) {
	teardownSuite := setupSuite()

	satellite, err := controller.CreateSatellite(withPermissions([]string{"satellite.create"}), controller.CreateSatelliteDto{
		Name:        "Test",
		Description: "description",
		Address:     "127.0.0.1:1234",
		Token:       "test",
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	satellite, err = controller.GetSatellite(withPermissions([]string{"satellite.view"}), controller.GetSatelliteDto{
		Id: satellite.Id,
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, "Test", satellite.Name)
	assert.Equal(t, "description", satellite.Description)
	assert.Equal(t, "127.0.0.1:1234", satellite.Address)

	teardownSuite()
}

func Test_GetSatellites(t *testing.T) {
	teardownSuite := setupSuite()

	_, err := controller.CreateSatellite(withPermissions([]string{"satellite.create"}), controller.CreateSatelliteDto{
		Name:        "Test",
		Description: "description",
		Address:     "127.0.0.1:1234",
		Token:       "test",
	})
	_, err = controller.CreateSatellite(withPermissions([]string{"satellite.create"}), controller.CreateSatelliteDto{
		Name:        "Test2",
		Description: "description2",
		Address:     "127.0.0.1:12314",
		Token:       "testasd",
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	satellites, err := controller.GetSatellites(withPermissions([]string{"satellite.view"}), controller.GetPageDto{})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}
	assert.Equal(t, 2, len(satellites.Rows))
	teardownSuite()
}

func Test_UpdateSatellite(t *testing.T) {
	teardownSuite := setupSuite()

	target, err := controller.CreateSatellite(withPermissions([]string{"satellite.create"}), controller.CreateSatelliteDto{
		Name:        "Test",
		Description: "description",
		Address:     "127.0.0.1:1234",
		Token:       "test",
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	updatedName := "Test2"
	satellite, err := controller.UpdateSatellite(withPermissions([]string{"satellite.update"}), controller.UpdateSatelliteDto{
		target.Id,
		&updatedName,
		nil,
		nil,
		nil,
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	satelliteResolve, err := controller.GetSatellite(withPermissions([]string{"satellite.view"}), controller.GetSatelliteDto{
		Id: satellite.Id,
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, updatedName, satelliteResolve.Name)
	teardownSuite()
}

func Test_DeleteSatellite(t *testing.T) {
	teardownSuite := setupSuite()

	target, err := controller.CreateSatellite(withPermissions([]string{"satellite.create"}), controller.CreateSatelliteDto{
		Name:        "Test",
		Description: "description",
		Address:     "127.0.0.1:1234",
		Token:       "test",
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	_, err = controller.DeleteSatellite(withPermissions([]string{"satellite.delete"}), controller.DeleteSatelliteDto{
		Id: target.Id,
	})

	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	satellites, err := controller.GetSatellites(withPermissions([]string{"satellite.view"}), controller.GetPageDto{})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, int64(0), satellites.TotalRows)
	teardownSuite()
}
