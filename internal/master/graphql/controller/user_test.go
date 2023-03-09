package controller_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"kroseida.org/slixx/internal/master/datasource"
	"kroseida.org/slixx/internal/master/graphql/controller"
	"testing"
)

func Test_CreateUser_MissingPermission(t *testing.T) {
	teardownSuite := setupSuite()
	_, err := controller.CreateUser(withPermissions([]string{"user.notcreate"}), controller.CreateUserDto{})
	if err == nil && err.Error() != "missing permission" {
		t.Error("Expected error (permission denied)")
		teardownSuite()
		return
	}
	teardownSuite()
}

func Test_CreateUser(t *testing.T) {
	teardownSuite := setupSuite()
	_, err := controller.CreateUser(withPermissions([]string{"user.create"}), controller.CreateUserDto{
		Name:        "Testaaaaaa",
		FirstName:   "test",
		LastName:    "test",
		Email:       "test@test.de",
		Description: "description",
		Active:      true,
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	actualUser, err := datasource.UserProvider.GetUserByName("Testaaaaaa")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.NotNil(t, actualUser)
	assert.Equal(t, "Testaaaaaa", actualUser.Name)
	assert.Equal(t, "test", actualUser.FirstName)
	assert.Equal(t, "test", actualUser.LastName)
	assert.Equal(t, "description", actualUser.Description)
	assert.Equal(t, "test@test.de", actualUser.Email)
	assert.Equal(t, true, actualUser.Active)
	teardownSuite()
}

func Test_DeleteUser_MissingPermission(t *testing.T) {
	teardownSuite := setupSuite()

	_, err := controller.DeleteUser(withPermissions([]string{"user.nodelete"}), controller.DeleteUserDto{})

	if err == nil && err.Error() != "missing permission" {
		t.Error("Expected error (permission denied)")
		teardownSuite()
		return
	}

	teardownSuite()
}

func Test_DeleteUser(t *testing.T) {
	teardownSuite := setupSuite()
	user, err := datasource.UserProvider.CreateUser(
		"oldName",
		"Test@test.de",
		"Test",
		"Test",
		"Test",
		true,
	)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	userInDatabase, err := datasource.UserProvider.GetUser(user.Id)
	assert.NotNil(t, userInDatabase)

	_, err = controller.DeleteUser(withPermissions([]string{"user.delete"}), controller.DeleteUserDto{
		Id: user.Id,
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	userInDatabase, err = datasource.UserProvider.GetUser(user.Id)
	assert.Nil(t, userInDatabase)

	teardownSuite()
}

func Test_AddUserPermission_MissingPermission(t *testing.T) {
	teardownSuite := setupSuite()
	_, err := controller.AddUserPermission(withPermissions([]string{"user.permission.notupdate"}), controller.AddUserPermissionDto{
		Id:          uuid.New(),
		Permissions: []string{"user.update"},
	})
	if err == nil && err.Error() != "missing permission" {
		t.Error("Expected error (permission denied)")
		teardownSuite()
		return
	}
	teardownSuite()
}

func Test_AddUserPermission(t *testing.T) {
	teardownSuite := setupSuite()

	user, err := datasource.UserProvider.CreateUser(
		"user",
		"Test@test.de",
		"Test",
		"Test",
		"Test",
		true,
	)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	_, err = controller.AddUserPermission(withPermissions([]string{"user.permission.update"}), controller.AddUserPermissionDto{
		Id:          user.Id,
		Permissions: []string{"user.update"},
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	actualUser, err := datasource.UserProvider.GetUserByName("user")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, "[\"user.update\"]", actualUser.Permissions)
	teardownSuite()
}

func Test_RemoveUserPermission_MissingPermission(t *testing.T) {
	teardownSuite := setupSuite()
	_, err := controller.RemoveUserPermission(withPermissions([]string{"user.permission.notupdate"}), controller.RemoveUserPermissionDto{
		Id:          uuid.New(),
		Permissions: []string{"user.update"},
	})
	if err == nil && err.Error() != "missing permission" {
		t.Error("Expected error (permission denied)")
		teardownSuite()
		return
	}
	teardownSuite()
}

func Test_RemoveUserPermission(t *testing.T) {
	teardownSuite := setupSuite()

	user, err := datasource.UserProvider.CreateUser(
		"user",
		"Test@test.de",
		"Test",
		"Test",
		"Test",
		true,
	)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	_, err = datasource.UserProvider.AddUserPermission(user.Id, []string{"user.update"})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	_, err = controller.RemoveUserPermission(withPermissions([]string{"user.permission.update"}), controller.RemoveUserPermissionDto{
		Id:          user.Id,
		Permissions: []string{"user.update"},
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	actualUser, err := datasource.UserProvider.GetUserByName("user")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, "[]", actualUser.Permissions)
	teardownSuite()
}

func Test_CreatePasswordAuthentication_MissingPermission(t *testing.T) {
	teardownSuite := setupSuite()
	_, err := controller.CreatePasswordAuthentication(withPermissions([]string{"user.notupdate"}), controller.UpdateUserPasswordDto{
		Id:       uuid.New(),
		Password: "123123123",
	})
	if err == nil && err.Error() != "missing permission" {
		t.Error("Expected error (permission denied)")
		teardownSuite()
		return
	}
	teardownSuite()
}

func Test_CreatePasswordAuthentication(t *testing.T) {
	teardownSuite := setupSuite()

	user, err := datasource.UserProvider.CreateUser(
		"user",
		"Test@test.de",
		"Test",
		"Test",
		"Test",
		true,
	)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	_, err = controller.CreatePasswordAuthentication(withPermissions([]string{"user.update"}), controller.UpdateUserPasswordDto{
		Id:       user.Id,
		Password: "123123123",
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	session, err := controller.Authenticate(context.Background(), controller.PasswordAuthenticationDto{
		Name:     "user",
		Password: "123123123",
	})

	assert.NotNil(t, session)
	teardownSuite()
}

func Test_GetUsers_MissingPermission(t *testing.T) {
	teardownSuite := setupSuite()
	_, err := controller.GetUsers(withPermissions([]string{"user.notview"}), controller.PageArgs{})
	if err == nil && err.Error() != "missing permission" {
		t.Error("Expected error (permission denied)")
		teardownSuite()
		return
	}
	teardownSuite()
}

func Test_GetUsers(t *testing.T) {
	teardownSuite := setupSuite()
	users, err := controller.GetUsers(withPermissions([]string{"user.view"}), controller.PageArgs{})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, 1, len(users.Rows))
	teardownSuite()
}

func Test_GetUser_MissingPermission(t *testing.T) {
	teardownSuite := setupSuite()
	_, err := controller.GetUser(withPermissions([]string{"user.notview"}), controller.GetUserDto{})
	if err == nil && err.Error() != "missing permission" {
		t.Error("Expected error (permission denied)")
		teardownSuite()
		return
	}
	teardownSuite()
}

func Test_GetUser(t *testing.T) {
	teardownSuite := setupSuite()

	userByName, err := datasource.UserProvider.GetUserByName("admin")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	user, err := controller.GetUser(withPermissions([]string{"user.view"}), controller.GetUserDto{
		Id: userByName.Id,
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, "admin", user.Name)
	teardownSuite()
}

func Test_GetLocalUser(t *testing.T) {
	teardownSuite := setupSuite()

	userByName, err := datasource.UserProvider.GetUserByName("admin")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	user, err := controller.GetLocalUser(context.WithValue(context.Background(), "user", userByName))
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, "admin", user.Name)
	teardownSuite()
}
