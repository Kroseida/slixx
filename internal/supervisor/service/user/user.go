package user

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/samsarahq/thunder/graphql"
	"kroseida.org/slixx/internal/supervisor/application"
	"kroseida.org/slixx/internal/supervisor/authenticator"
	"kroseida.org/slixx/internal/supervisor/datasource"
	"kroseida.org/slixx/internal/supervisor/datasource/provider"
	"kroseida.org/slixx/pkg/model"
	"time"
)

func AddPermission(userId uuid.UUID, permissions []string) (*model.User, error) {
	return datasource.UserProvider.AddPermission(userId, permissions)
}

func RemovePermission(userId uuid.UUID, permissions []string) (*model.User, error) {
	return datasource.UserProvider.RemovePermission(userId, permissions)
}

func CreatePasswordAuthentication(userId uuid.UUID, password string) (*model.Authentication, error) {
	user, err := datasource.UserProvider.Get(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, graphql.NewSafeError("user not found")
	}

	// Remove any existing password authentications for this user
	datasource.UserProvider.DeleteAuthenticationOfKind(authenticator.PASSWORD, userId)

	kind := authenticator.GetKind(authenticator.PASSWORD).(authenticator.Password)

	configuration, err := kind.GenerateConfigurationFromRequestContainer(&authenticator.PasswordRequestContainer{
		Name:     user.Name,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	configurationJson, err := json.Marshal(configuration)
	if err != nil {
		return nil, err
	}

	return datasource.UserProvider.CreateAuthentication(userId, authenticator.PASSWORD, string(configurationJson))
}

func Authenticate(kindName string, configuration string) (*model.Session, error) {
	user, err := datasource.UserProvider.ValidateByAuthentication(kindName, configuration)
	if err != nil {
		return nil, err
	}

	return datasource.UserProvider.CreateSession(
		user.Id,
		time.Now().Add(time.Hour*time.Duration(application.CurrentSettings.Authentication.SessionDuration)),
	)
}

func AuthenticatePassword(name string, password string) (*model.Session, error) {
	configuration, err := json.Marshal(authenticator.PasswordRequestContainer{
		Name:     name,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	return Authenticate(authenticator.PASSWORD, string(configuration))
}

func GetPaged(pagination *provider.Pagination[model.User]) (*provider.Pagination[model.User], error) {
	return datasource.UserProvider.ListPaged(pagination)
}

func Get(id uuid.UUID) (*model.User, error) {
	return datasource.UserProvider.Get(id)
}

func Create(
	name string,
	email string,
	firstName string,
	lastName string,
	description string,
	active bool,
) (*model.User, error) {
	return datasource.UserProvider.Create(name, email, firstName, lastName, description, active)
}

func Update(
	id uuid.UUID,
	name *string,
	firstName *string,
	lastName *string,
	active *bool,
	description *string,
	email *string,
) (*model.User, error) {
	return datasource.UserProvider.Update(id, name, firstName, lastName, active, description, email)
}

func GetUserBySession(token string) (uuid.UUID, error) {
	return datasource.UserProvider.GetBySession(token)
}

func Delete(id uuid.UUID) (*model.User, error) {
	return datasource.UserProvider.Delete(id)
}
