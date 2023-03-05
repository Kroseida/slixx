package provider_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"kroseida.org/slixx/internal/master/datasource"
	"testing"
	"time"
)

func Test_CreateUser(t *testing.T) {
	teardownSuite := setupSuite()

	_, err := datasource.UserProvider.CreateUser(
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

	user, err := datasource.UserProvider.GetUserByName("Test")
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

	_, err := datasource.UserProvider.CreateUser(
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

	_, err := datasource.UserProvider.CreateUser(
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

	_, err := datasource.UserProvider.CreateUser(
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

	_, err = datasource.UserProvider.CreateUser(
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

func Test_DeleteUser(t *testing.T) {
	teardownSuite := setupSuite()

	user, err := datasource.UserProvider.CreateUser(
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

	users, err := datasource.UserProvider.GetUsers()
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, 2, len(users))

	_, err = datasource.UserProvider.DeleteUser(user.Id)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	users, err = datasource.UserProvider.GetUsers()
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, 1, len(users))

	teardownSuite()
}

func Test_GetUsers(t *testing.T) {
	teardownSuite := setupSuite()

	_, err := datasource.UserProvider.CreateUser(
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

	_, err = datasource.UserProvider.CreateUser(
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

	users, err := datasource.UserProvider.GetUsers()
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

	_, err := datasource.UserProvider.CreateUser(
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

	createdUser, err := datasource.UserProvider.CreateUser(
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

	user, err := datasource.UserProvider.GetUser(createdUser.Id)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, "Max", user.Name)
	teardownSuite()
}

func Test_CreateAuthentication_MissingUser(t *testing.T) {
	teardownSuite := setupSuite()

	userId := uuid.New()

	_, err := datasource.UserProvider.CreatePasswordAuthentication(userId, "password")
	if err == nil {
		t.Error("Expected error (missing user)")
		teardownSuite()
		return
	}
	teardownSuite()
}

func Test_GetUserBySession(t *testing.T) {
	teardownSuite := setupSuite()

	createdUser, err := datasource.UserProvider.CreateUser(
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
	userId, err := datasource.UserProvider.GetUserBySession(session.Token)
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

	createdUser, err := datasource.UserProvider.CreateUser(
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
	userId, err := datasource.UserProvider.GetUserBySession(session.Token + "_invalid")
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

	createdUser, err := datasource.UserProvider.CreateUser(
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
	sessions, err := datasource.UserProvider.GetSessions(createdUser.Id)
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

	createdUser, err := datasource.UserProvider.CreateUser(
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
	sessions, err := datasource.UserProvider.GetSessions(createdUser.Id)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}
	assert.Equal(t, 0, len(sessions))
	teardownSuite()
}

func Test_Authenticate(t *testing.T) {
	teardownSuite := setupSuite()

	createdUser, err := datasource.UserProvider.CreateUser(
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
	password := "1234!Test*_:"
	_, err = datasource.UserProvider.CreatePasswordAuthentication(createdUser.Id, password)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}
	session, err := datasource.UserProvider.AuthenticatePassword("Max", password)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, createdUser.Id, session.UserId)
	teardownSuite()
}

func Test_Authenticate_Invalid(t *testing.T) {
	teardownSuite := setupSuite()

	createdUser, err := datasource.UserProvider.CreateUser(
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
	password := "1234!Test*_:"
	_, err = datasource.UserProvider.CreatePasswordAuthentication(createdUser.Id, password)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}
	_, err = datasource.UserProvider.AuthenticatePassword("Max", "123123123")
	if err == nil {
		t.Error("Expected error (invalid password)")
		teardownSuite()
		return
	}
	teardownSuite()
}

func Test_Authenticate_InvalidUser(t *testing.T) {
	teardownSuite := setupSuite()

	createdUser, err := datasource.UserProvider.CreateUser(
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
	password := "1234!Test*_:"
	_, err = datasource.UserProvider.CreatePasswordAuthentication(createdUser.Id, password)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}
	_, err = datasource.UserProvider.AuthenticatePassword("Alex", password)
	if err == nil {
		t.Error("Expected error (invalid password)")
		teardownSuite()
		return
	}
	teardownSuite()
}
