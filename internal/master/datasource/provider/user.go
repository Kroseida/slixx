package provider

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/samsarahq/thunder/graphql"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"kroseida.org/slixx/internal/master/application"
	"kroseida.org/slixx/internal/master/authenticator"
	"kroseida.org/slixx/internal/master/datasource/model"
	"kroseida.org/slixx/pkg/utils"
	"strings"
	"time"
)

// UserProvider User Provider
type UserProvider struct {
	Database *gorm.DB
}

func (provider UserProvider) CreateUser(name string, email string, firstName string, lastName string, description string, active bool) (*model.User, error) {
	if name == "" {
		return nil, graphql.NewSafeError("name can not be empty")
	}
	if strings.Contains(name, " ") {
		return nil, graphql.NewSafeError("name can not contain spaces")
	}

	existingUser, err := provider.GetUserByName(name)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, graphql.NewSafeError("name already in use")
	}

	user := &model.User{
		Id:          uuid.New(),
		Name:        name,
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		Description: description,
		Active:      active,
		Permissions: "[]",
	}

	result := provider.Database.Create(&user)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return user, nil
}

func (provider UserProvider) UpdateUserName(id uuid.UUID, name string) (*model.User, error) {
	user, err := provider.GetUser(id)
	if user == nil {
		return nil, graphql.NewSafeError("user not found")
	}
	if err != nil {
		return nil, err
	}

	if strings.Contains(name, " ") {
		return nil, graphql.NewSafeError("name can not contain spaces")
	}
	existingUser, err := provider.GetUserByName(name)
	if existingUser != nil {
		return nil, graphql.NewSafeError("name already in use")
	}
	if err != nil {
		return nil, err
	}

	user.Name = name

	result := provider.Database.Save(&user)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return user, nil
}

func (provider UserProvider) UpdateUserFirstName(id uuid.UUID, firstName string) (*model.User, error) {
	user, err := provider.GetUser(id)
	if user == nil {
		return nil, graphql.NewSafeError("user not found")
	}
	if err != nil {
		return nil, err
	}

	user.FirstName = firstName

	result := provider.Database.Save(&user)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return user, nil
}

func (provider UserProvider) UpdateUserLastName(id uuid.UUID, lastName string) (*model.User, error) {
	user, err := provider.GetUser(id)
	if user == nil {
		return nil, graphql.NewSafeError("user not found")
	}
	if err != nil {
		return nil, err
	}

	user.LastName = lastName

	result := provider.Database.Save(&user)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return user, nil
}

func (provider UserProvider) UpdateUserEmail(id uuid.UUID, email string) (*model.User, error) {
	user, err := provider.GetUser(id)
	if user == nil {
		return nil, graphql.NewSafeError("user not found")
	}
	if err != nil {
		return nil, err
	}

	user.Email = email

	result := provider.Database.Save(&user)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return user, nil
}

func (provider UserProvider) UpdateUserActive(id uuid.UUID, active bool) (*model.User, error) {
	user, err := provider.GetUser(id)
	if user == nil {
		return nil, graphql.NewSafeError("user not found")
	}
	if err != nil {
		return nil, err
	}

	user.Active = active

	result := provider.Database.Save(&user)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return user, nil
}

func (provider UserProvider) UpdateUserDescription(id uuid.UUID, description string) (*model.User, error) {
	user, err := provider.GetUser(id)
	if user == nil {
		return nil, graphql.NewSafeError("user not found")
	}
	if err != nil {
		return nil, err
	}

	user.Description = description

	result := provider.Database.Save(&user)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return user, nil
}

func (provider UserProvider) AddUserPermission(userId uuid.UUID, permissions []string) (*model.User, error) {
	user, err := provider.GetUser(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, graphql.NewSafeError("user not found")
	}

	user.AddPermission(permissions)

	result := provider.Database.Save(&user)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return user, nil
}

func (provider UserProvider) RemoveUserPermission(userId uuid.UUID, permissions []string) (*model.User, error) {
	user, err := provider.GetUser(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, graphql.NewSafeError("user not found")
	}

	user.RemovePermission(permissions)

	result := provider.Database.Save(&user)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return user, nil
}

func (provider UserProvider) GetUsers() ([]*model.User, error) {
	var users []*model.User
	result := provider.Database.Find(&users)

	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return users, nil
}

func (provider UserProvider) GetUser(id uuid.UUID) (*model.User, error) {
	var user *model.User
	result := provider.Database.First(&user, "id = ?", id)

	if isSqlError(result.Error) {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return user, nil
}

func (provider UserProvider) GetUserByName(name string) (*model.User, error) {
	var user *model.User
	result := provider.Database.First(&user, "name = ?", name)

	if isSqlError(result.Error) {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return user, nil
}

func (provider UserProvider) CreateSession(userId uuid.UUID, expiresAt time.Time) (*model.Session, error) {
	user, err := provider.GetUser(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, graphql.NewSafeError("user not found")
	}

	session := &model.Session{
		Id:        uuid.New(),
		UserId:    userId,
		ExpiresAt: expiresAt,
		Token:     utils.GenerateSecureToken(application.CurrentSettings.Authentication.TokenSize),
	}

	result := provider.Database.Create(&session)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return session, nil
}

func (provider UserProvider) GetSessions(userId uuid.UUID) ([]*model.Session, error) {
	user, err := provider.GetUser(userId)
	if err != nil {
		return []*model.Session{}, err
	}
	if user == nil {
		return []*model.Session{}, graphql.NewSafeError("user not found")
	}

	var sessions []*model.Session
	result := provider.Database.Find(&sessions, "user_id = ? AND expires_at > ?", userId, time.Now())

	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return sessions, nil
}

func (provider UserProvider) GetUserBySession(token string) (uuid.UUID, error) {
	var session *model.Session
	result := provider.Database.First(&session, "token = ? AND expires_at > ?", token, time.Now())

	if isSqlError(result.Error) {
		return uuid.UUID{}, result.Error
	}
	if result.RowsAffected == 0 {
		return uuid.UUID{}, nil
	}

	return session.UserId, nil
}

func (provider UserProvider) AuthenticatePassword(name string, password string) (*model.Session, error) {
	configuration, err := json.Marshal(authenticator.PasswordRequestContainer{
		Name:     name,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	return provider.Authenticate(authenticator.PASSWORD, string(configuration))
}

func (provider UserProvider) Authenticate(kindName string, configuration string) (*model.Session, error) {
	kind := authenticator.GetKind(kindName)
	if kind == nil {
		return nil, graphql.NewSafeError("invalid authentication kind")
	}

	var authentications []*model.Authentication
	result := provider.Database.Find(&authentications, "kind = ?", kindName)

	if isSqlError(result.Error) {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	hasBcryptError := false

	var targetAuthentication *model.Authentication
	for _, authentication := range authentications {
		isValid, err := kind.Validate(authentication, configuration)
		if err != nil {
			if err == bcrypt.ErrMismatchedHashAndPassword {
				hasBcryptError = true
			}
			continue
		}
		if isValid {
			targetAuthentication = authentication
			break
		}
	}
	if targetAuthentication == nil {
		if !hasBcryptError {
			// we don't want to leak information about the existence of a user with the given name over timing attacks
			// (auth with name a: took 500ms -> user does not exist, because system did not ran into bcrypt hash)
			// so we delay the response based on the bcrypt hashing duration.
			// (todo: make this configurable and improve for more security)
			bcrypt.GenerateFromPassword([]byte(utils.GenerateSecureToken(16)), application.CurrentSettings.Authentication.HashCost)
		}
		return nil, graphql.NewSafeError("authentication failed")
	}

	return provider.CreateSession(targetAuthentication.UserId, time.Now().Add(time.Hour*time.Duration(application.CurrentSettings.Authentication.SessionDuration)))
}

func (provider UserProvider) CreatePasswordAuthentication(userId uuid.UUID, password string) (*model.Authentication, error) {
	user, err := provider.GetUser(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, graphql.NewSafeError("user not found")
	}
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

	var authentication *model.Authentication
	result := provider.Database.Delete(&authentication, "kind = ? AND user_id = ?", authenticator.PASSWORD, userId)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return provider.CreateAuthentication(userId, authenticator.PASSWORD, string(configurationJson))
}

func (provider UserProvider) CreateAuthentication(userId uuid.UUID, kind string, configuration string) (*model.Authentication, error) {
	user, err := provider.GetUser(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, graphql.NewSafeError("user not found")
	}

	authentication := &model.Authentication{
		Id:            uuid.New(),
		Kind:          kind,
		Configuration: configuration,
		UserId:        userId,
	}

	_, err = authenticator.GetKind(kind).ParseConfiguration(configuration)
	if err != nil {
		return nil, err
	}

	result := provider.Database.Create(&authentication)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return authentication, nil
}

func (provider UserProvider) Init() error {
	// Create default user if migration did not create one
	err := provider.defaultUserMigration()
	if err != nil {
		return err
	}

	return nil
}

func (provider UserProvider) defaultUserMigration() error {
	var migration *model.Migration
	result := provider.Database.First(&migration, "name = ?", "default_user")
	if isSqlError(result.Error) {
		return result.Error
	}
	if result.RowsAffected != 0 {
		return nil
	}

	user, err := provider.CreateUser("admin", "", "", "", "default admin user", true)
	if err != nil {
		return err
	}
	_, err = provider.AddUserPermission(user.Id, []string{
		"user.create",
		"user.view",
		"user.update",
		"user.permission.update",
	})
	if err != nil {
		return err
	}

	_, err = provider.CreatePasswordAuthentication(user.Id, "admin")
	if err != nil {
		return err
	}

	err = provider.Database.Create(&model.Migration{
		Id:   uuid.New(),
		Name: "default_user",
	}).Error
	if err != nil {
		return err
	}

	return nil
}
