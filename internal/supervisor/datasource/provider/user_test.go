package provider_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"kroseida.org/slixx/internal/supervisor/datasource"
	"kroseida.org/slixx/internal/supervisor/datasource/provider"
	"kroseida.org/slixx/pkg/model"
	"testing"
	"time"
)

func Test_CreateUser(t *testing.T) {
	teardownSuite := setupSuite()

	_, err := datasource.UserProvider.Create(
		"Test",
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

	user, err := datasource.UserProvider.GetByName("Test")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, "Test", user.Name)
	teardownSuite()
}

func Test_CreateUser_EmptyName(t *testing.T) {
	teardownSuite := setupSuite()

	_, err := datasource.UserProvider.Create(
		"",
		"Test@test.de",
		"Test",
		"Test",
		"Test",
		true,
	)
	if err == nil {
		t.Error("Expected error (empty name)")
		teardownSuite()
		return
	}
	teardownSuite()
}

func Test_CreateUser_SpacedName(t *testing.T) {
	teardownSuite := setupSuite()

	_, err := datasource.UserProvider.Create(
		"Test User",
		"Test@test.de",
		"Test",
		"Test",
		"Test",
		true,
	)
	if err == nil {
		t.Error("Expected error (spaced name)")
		teardownSuite()
		return
	}
	teardownSuite()
}

func Test_CreateUser_NameUsed(t *testing.T) {
	teardownSuite := setupSuite()

	_, err := datasource.UserProvider.Create(
		"Test",
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

	_, err = datasource.UserProvider.Create(
		"Test",
		"Test@test.de",
		"Test",
		"Test",
		"Test",
		true,
	)
	if err == nil {
		t.Error("Expected error (name is use)")
		teardownSuite()
		return
	}
	teardownSuite()
}

func Test_UpdateUser_Email(t *testing.T) {
	teardownSuite := setupSuite()

	newEmail := "test"
	user, err := datasource.UserProvider.Create(
		"Test",
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
	_, err = datasource.UserProvider.Update(user.Id, nil, nil, nil, nil, nil, &newEmail)
	assert.NotNil(t, err)

	newEmail = "test@test.de"
	_, err = datasource.UserProvider.Update(user.Id, nil, nil, nil, nil, nil, &newEmail)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}
	updatedUser, err := datasource.UserProvider.Get(user.Id)

	assert.Equal(t, newEmail, updatedUser.Email)

	teardownSuite()
}

func Test_UpdateUser_Name(t *testing.T) {
	teardownSuite := setupSuite()

	newName := "Test2"
	user, err := datasource.UserProvider.Create(
		"Test",
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

	_, err = datasource.UserProvider.Update(user.Id, &newName, nil, nil, nil, nil, nil)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}
	updatedUser, err := datasource.UserProvider.Get(user.Id)

	assert.Equal(t, newName, updatedUser.Name)

	teardownSuite()
}

func Test_DeleteUser(t *testing.T) {
	teardownSuite := setupSuite()

	user, err := datasource.UserProvider.Create(
		"Test",
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

	users, err := datasource.UserProvider.List()
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, 2, len(users))

	_, err = datasource.UserProvider.Delete(user.Id)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	users, err = datasource.UserProvider.List()
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, 1, len(users))

	teardownSuite()
}

func Test_GetUsersPaged(t *testing.T) {
	teardownSuite := setupSuite()

	_, err := datasource.UserProvider.Create(
		"Test",
		"Test@test.de",
		"Test",
		"Test",
		"Test",
		true,
	)
	_, err = datasource.UserProvider.Create(
		"Test2",
		"Test2@test.de",
		"Test2",
		"Test2",
		"Test2",
		true,
	)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	users, err := datasource.UserProvider.ListPaged(&provider.Pagination[model.User]{
		Page:  1,
		Limit: 2,
	})

	assert.Equal(t, 2, len(users.Rows))
	teardownSuite()
}

func Test_GetUsers(t *testing.T) {
	teardownSuite := setupSuite()

	_, err := datasource.UserProvider.Create(
		"Test",
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

	_, err = datasource.UserProvider.Create(
		"Test2",
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

	users, err := datasource.UserProvider.List()
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, 3, len(users))
	if users[0].Name != "Test" && users[0].Name != "Test2" && users[0].Name != "admin" {
		t.Error("Expected user name 'Test' or 'Test2'")
	}
	if users[1].Name != "Test" && users[1].Name != "Test2" && users[1].Name != "admin" {
		t.Error("Expected user name 'Test' or 'Test2'")
	}
	if users[0].Name == users[1].Name || users[0].Name == users[2].Name || users[1].Name == users[2].Name {
		t.Error("Expected different user names")
	}
	teardownSuite()
}

func Test_GetUser(t *testing.T) {
	teardownSuite := setupSuite()

	_, err := datasource.UserProvider.Create(
		"Maxi",
		"Maxi@test.de",
		"Maxi",
		"Test",
		"Test",
		true,
	)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	createdUser, err := datasource.UserProvider.Create(
		"Max",
		"Max@test.de",
		"Max",
		"Test",
		"Test",
		true,
	)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	user, err := datasource.UserProvider.Get(createdUser.Id)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, "Max", user.Name)
	teardownSuite()
}

func Test_GetUserBySession(t *testing.T) {
	teardownSuite := setupSuite()

	createdUser, err := datasource.UserProvider.Create(
		"Max",
		"Max@test.de",
		"Max",
		"Test",
		"Test",
		true,
	)

	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}
	session, err := datasource.UserProvider.CreateSession(createdUser.Id, time.Now().Add(time.Hour))
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}
	userId, err := datasource.UserProvider.GetBySession(session.Token)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}
	assert.Equal(t, createdUser.Id, userId)
	teardownSuite()
}

func Test_GetUserBySession_Invalid(t *testing.T) {
	teardownSuite := setupSuite()

	createdUser, err := datasource.UserProvider.Create(
		"Max",
		"Max@test.de",
		"Max",
		"Test",
		"Test",
		true,
	)

	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}
	session, err := datasource.UserProvider.CreateSession(createdUser.Id, time.Now().Add(time.Hour))
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}
	userId, err := datasource.UserProvider.GetBySession(session.Token + "_invalid")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}
	assert.Equal(t, uuid.UUID{}, userId)
	teardownSuite()
}

func Test_CreateSession(t *testing.T) {
	teardownSuite := setupSuite()

	createdUser, err := datasource.UserProvider.Create(
		"Max",
		"Max@test.de",
		"Max",
		"Test",
		"Test",
		true,
	)

	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}
	_, err = datasource.UserProvider.CreateSession(createdUser.Id, time.Now().Add(time.Hour))
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}
	sessions, err := datasource.UserProvider.ListSessions(createdUser.Id)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}
	assert.Equal(t, 1, len(sessions))
	teardownSuite()
}

func Test_CreateSession_MissingUser(t *testing.T) {
	teardownSuite := setupSuite()

	userId := uuid.New()

	_, err := datasource.UserProvider.CreateSession(userId, time.Now().Add(time.Hour))
	if err == nil {
		t.Error("Expected error (missing user)")
		teardownSuite()
		return
	}
	teardownSuite()
}

func Test_GetSession_Expired(t *testing.T) {
	teardownSuite := setupSuite()

	createdUser, err := datasource.UserProvider.Create(
		"Max",
		"Max@test.de",
		"Max",
		"Test",
		"Test",
		true,
	)

	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}
	_, err = datasource.UserProvider.CreateSession(createdUser.Id, time.Now().Add(-time.Hour))
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}
	sessions, err := datasource.UserProvider.ListSessions(createdUser.Id)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}
	assert.Equal(t, 0, len(sessions))
	teardownSuite()
}

func Test_AddUserPermission(t *testing.T) {
	teardownSuite := setupSuite()

	user, err := datasource.UserProvider.Create(
		"Test",
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

	_, err = datasource.UserProvider.AddPermission(user.Id, []string{"test"})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	resolvedUser, err := datasource.UserProvider.Get(user.Id)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, "[\"test\"]", resolvedUser.Permissions)
	teardownSuite()
}

func Test_RemoveUserPermission(t *testing.T) {
	teardownSuite := setupSuite()

	user, err := datasource.UserProvider.Create(
		"Test",
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
	_, err = datasource.UserProvider.AddPermission(user.Id, []string{"test", "test2", "test3"})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	_, err = datasource.UserProvider.RemovePermission(user.Id, []string{"test2", "test3"})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	resolvedUser, err := datasource.UserProvider.Get(user.Id)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, "[\"test\"]", resolvedUser.Permissions)
	teardownSuite()
}
