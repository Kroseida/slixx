package authenticator_test

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"kroseida.org/slixx/internal/master/authenticator"
	"kroseida.org/slixx/pkg/model"
	"testing"
)

func Test_Validate(t *testing.T) {
	t.Parallel()
	kind := authenticator.GetKind(authenticator.PASSWORD).(authenticator.Password)
	authentication, err := setupAuthentication("test", "123123123")
	if err != nil {
		t.Error(err)
		return
	}

	request := authenticator.PasswordRequestContainer{
		Name:     "test",
		Password: "123123123",
	}
	requestJson, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
		return
	}

	isValid, err := kind.Validate(authentication, string(requestJson))

	assert.Equal(t, true, isValid)
}

func Test_Validate_WrongHash(t *testing.T) {
	t.Parallel()
	kind := authenticator.GetKind(authenticator.PASSWORD).(authenticator.Password)

	configuration := map[string]string{
		"name": "test",
		"hash": "invalidHash",
	}

	configurationJson, err := json.Marshal(configuration)
	if err != nil {
		t.Error(err)
		return
	}

	authentication := model.Authentication{
		Id:            uuid.New(),
		Kind:          authenticator.PASSWORD,
		Configuration: string(configurationJson),
		UserId:        uuid.New(),
	}

	if err != nil {
		t.Error(err)
		return
	}

	request := authenticator.PasswordRequestContainer{
		Name:     "test",
		Password: "123123123",
	}
	requestJson, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = kind.Validate(&authentication, string(requestJson))
	if err == nil {
		t.Fail()
		return
	}
}

func Test_Validate_Wrong_Name(t *testing.T) {
	t.Parallel()
	kind := authenticator.GetKind(authenticator.PASSWORD).(authenticator.Password)
	authentication, err := setupAuthentication("test", "123123123")
	if err != nil {
		t.Error(err)
		return
	}

	request := authenticator.PasswordRequestContainer{
		Name:     "Max",
		Password: "123123123",
	}
	requestJson, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
		return
	}

	isValid, err := kind.Validate(authentication, string(requestJson))

	assert.Equal(t, false, isValid)
}

func Test_Validate_Wrong_Password(t *testing.T) {
	kind := authenticator.GetKind(authenticator.PASSWORD).(authenticator.Password)
	authentication, err := setupAuthentication("test", "123123123")
	if err != nil {
		t.Error(err)
		return
	}

	request := authenticator.PasswordRequestContainer{
		Name:     "test",
		Password: "test",
	}
	requestJson, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
		return
	}

	isValid, err := kind.Validate(authentication, string(requestJson))

	assert.Equal(t, false, isValid)
}

func Test_ParseRequestContainer(t *testing.T) {
	t.Parallel()
	kind := authenticator.GetKind(authenticator.PASSWORD).(authenticator.Password)

	passwordConfiguration := authenticator.PasswordRequestContainer{
		Name:     "test",
		Password: "123123123",
	}
	passwordConfigurationJson, err := json.Marshal(passwordConfiguration)
	if err != nil {
		t.Error(err)
		return
	}
	request, err := kind.ParseRequestContainer(string(passwordConfigurationJson))

	assert.Equal(t, passwordConfiguration.Name, request.Name)
	assert.Equal(t, passwordConfiguration.Password, request.Password)
}

func Test_ParseRequestContainer_MissingName(t *testing.T) {
	kind := authenticator.GetKind(authenticator.PASSWORD).(authenticator.Password)

	passwordConfiguration := map[string]string{
		"password": "test",
	}
	passwordConfigurationJson, err := json.Marshal(passwordConfiguration)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = kind.ParseRequestContainer(string(passwordConfigurationJson))
	if err == nil {
		t.Fail()
		return
	}
}

func Test_ParseRequestContainer_MissingPassword(t *testing.T) {
	t.Parallel()
	kind := authenticator.GetKind(authenticator.PASSWORD).(authenticator.Password)

	passwordConfiguration := map[string]string{
		"name": "test",
	}
	passwordConfigurationJson, err := json.Marshal(passwordConfiguration)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = kind.ParseRequestContainer(string(passwordConfigurationJson))
	if err == nil {
		t.Fail()
		return
	}
}

func Test_ParseConfiguration(t *testing.T) {
	t.Parallel()
	kind := authenticator.GetKind(authenticator.PASSWORD).(authenticator.Password)

	passwordConfiguration := authenticator.PasswordConfiguration{
		Name: "test",
		Hash: "123123123",
	}
	passwordConfigurationJson, err := json.Marshal(passwordConfiguration)
	if err != nil {
		t.Error(err)
		return
	}

	configuration, err := kind.ParseConfiguration(string(passwordConfigurationJson))
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, passwordConfiguration.Name, configuration.(*authenticator.PasswordConfiguration).Name)
	assert.Equal(t, passwordConfiguration.Hash, configuration.(*authenticator.PasswordConfiguration).Hash)
}

func Test_ParseConfiguration_Without_Hash(t *testing.T) {
	t.Parallel()
	kind := authenticator.GetKind(authenticator.PASSWORD).(authenticator.Password)

	passwordConfiguration := map[string]string{
		"name": "test",
	}

	passwordConfigurationJson, err := json.Marshal(passwordConfiguration)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = kind.ParseConfiguration(string(passwordConfigurationJson))
	if err == nil {
		t.Fail()
		return
	}
}

func Test_ParseConfiguration_Without_Name(t *testing.T) {
	t.Parallel()
	kind := authenticator.GetKind(authenticator.PASSWORD).(authenticator.Password)

	passwordConfiguration := map[string]string{
		"hash": "123123123",
	}

	passwordConfigurationJson, err := json.Marshal(passwordConfiguration)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = kind.ParseConfiguration(string(passwordConfigurationJson))
	if err == nil {
		t.Fail()
		return
	}
}

func Test_GenerateConfigurationFromRequestContainer(t *testing.T) {
	t.Parallel()
	kind := authenticator.GetKind(authenticator.PASSWORD).(authenticator.Password)

	configuration, err := kind.GenerateConfigurationFromRequestContainer(&authenticator.PasswordRequestContainer{
		Name:     "test",
		Password: "123123123",
	})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, "test", configuration.Name)
	assert.NotEqual(t, "123123123", configuration.Hash)
}

func setupAuthentication(name string, password string) (*model.Authentication, error) {
	kind := authenticator.GetKind(authenticator.PASSWORD).(authenticator.Password)

	configuration, err := kind.GenerateConfigurationFromRequestContainer(&authenticator.PasswordRequestContainer{
		Name:     name,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	configurationJson, err := json.Marshal(configuration)
	if err != nil {
		return nil, err
	}

	authentication := model.Authentication{
		Id:            uuid.New(),
		Kind:          authenticator.PASSWORD,
		Configuration: string(configurationJson),
		UserId:        uuid.New(),
	}
	return &authentication, nil
}
