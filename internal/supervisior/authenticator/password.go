package authenticator

import (
	"encoding/json"
	"github.com/samsarahq/thunder/graphql"
	"golang.org/x/crypto/bcrypt"
	"kroseida.org/slixx/internal/supervisior/application"
	"kroseida.org/slixx/pkg/model"
)

type Password struct {
}

type PasswordConfiguration struct {
	Name string
	Hash string
}

type PasswordRequestContainer struct {
	Name     string
	Password string
}

func (kind Password) Validate(authentication *model.Authentication, sendConfigurationJson string) (bool, error) {
	sendConfiguration, err := kind.ParseRequestContainer(sendConfigurationJson)
	if err != nil {
		return false, err
	}

	configurationRaw, err := kind.ParseConfiguration(authentication.Configuration)
	if err != nil {
		return false, err
	}

	var configuration = configurationRaw.(*PasswordConfiguration)

	if configuration.Name != sendConfiguration.Name {
		return false, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(configuration.Hash), []byte(sendConfiguration.Password))

	if err != nil {
		return false, err
	}

	return true, nil
}

func (Password) ParseRequestContainer(configurationJson string) (*PasswordRequestContainer, error) {
	var container PasswordRequestContainer
	err := json.Unmarshal([]byte(configurationJson), &container)
	if err != nil {
		return nil, err
	}

	if container.Name == "" {
		return nil, graphql.NewSafeError("name is required")
	}
	if container.Password == "" {
		return nil, graphql.NewSafeError("password is required")
	}

	return &container, nil
}

func (Password) ParseConfiguration(configurationJson string) (any, error) {
	var configuration PasswordConfiguration
	err := json.Unmarshal([]byte(configurationJson), &configuration)
	if err != nil {
		return nil, err
	}

	if configuration.Name == "" {
		return nil, graphql.NewSafeError("name is required")
	}
	if configuration.Hash == "" {
		return nil, graphql.NewSafeError("hash is required")
	}

	return &configuration, nil
}

func (Password) GenerateConfigurationFromRequestContainer(configuration *PasswordRequestContainer) (*PasswordConfiguration, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(configuration.Password), application.CurrentSettings.Authentication.HashCost)

	if err != nil {
		return nil, err
	}
	return &PasswordConfiguration{
		configuration.Name,
		string(hash),
	}, nil
}
