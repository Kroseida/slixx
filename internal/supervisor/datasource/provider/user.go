package provider

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/samsarahq/thunder/graphql"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"kroseida.org/slixx/internal/supervisor/application"
	"kroseida.org/slixx/internal/supervisor/authenticator"
	"kroseida.org/slixx/pkg/model"
	"kroseida.org/slixx/pkg/utils"
	"strings"
	"time"
)

// UserProvider User Provider
type UserProvider struct {
	Database *gorm.DB
}

var Permissions = map[string]string{
	"user.view":              "View User",
	"user.create":            "Create User",
	"user.update":            "Update User Account",
	"user.delete":            "Delete User",
	"user.permission.update": "Update User Permissions",
	"storage.view":           "View Storage",
	"storage.create":         "Create Storage",
	"storage.update":         "Update Storage",
	"storage.delete":         "Delete Storage",
	"job.view":               "View Job",
	"job.create":             "Create Job",
	"job.update":             "Update Job",
	"job.delete":             "Delete Job",
	"satellite.view":         "View Satellite",
	"satellite.create":       "Register Satellite",
	"satellite.update":       "Update Satellite",
	"satellite.delete":       "Delete Satellite",
}

func (provider UserProvider) Create(name string, email string, firstName string, lastName string, description string, active bool) (*model.User, error) {
	if name == "" {
		return nil, graphql.NewSafeError("name can not be empty")
	}
	if email != "" && !strings.Contains(email, "@") {
		return nil, graphql.NewSafeError("invalid email")
	}
	if strings.Contains(name, " ") {
		return nil, graphql.NewSafeError("name can not contain spaces")
	}

	existingUser, err := provider.GetByName(name)
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

func (provider UserProvider) Delete(id uuid.UUID) (*model.User, error) {
	user, err := provider.Get(id)
	if user == nil {
		return nil, graphql.NewSafeError("user not found")
	}
	if err != nil {
		return nil, err
	}

	result := provider.Database.Delete(&user)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return user, nil
}

func (provider UserProvider) Update(
	id uuid.UUID,
	name *string,
	firstName *string,
	lastName *string,
	active *bool,
	description *string,
	email *string,
) (*model.User, error) {
	user, err := provider.Get(id)
	if user == nil {
		return nil, graphql.NewSafeError("user not found")
	}
	if err != nil {
		return nil, err
	}

	if name != nil {
		if strings.Contains(*name, " ") {
			return nil, graphql.NewSafeError("name can not contain spaces")
		}
		existingUser, err := provider.GetByName(*name)
		if existingUser != nil {
			return nil, graphql.NewSafeError("name already in use")
		}
		if err != nil {
			return nil, err
		}
		user.Name = *name
	}
	if firstName != nil {
		user.FirstName = *firstName
	}
	if lastName != nil {
		user.LastName = *lastName
	}
	if active != nil {
		user.Active = *active
	}
	if description != nil {
		user.Description = *description
	}
	if email != nil {
		if *email != "" && !strings.Contains(*email, "@") {
			return nil, graphql.NewSafeError("invalid email")
		}
		user.Email = *email
	}

	result := provider.Database.Save(&user)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return user, nil
}

func (provider UserProvider) AddPermission(userId uuid.UUID, permissions []string) (*model.User, error) {
	user, err := provider.Get(userId)
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

func (provider UserProvider) RemovePermission(userId uuid.UUID, permissions []string) (*model.User, error) {
	user, err := provider.Get(userId)
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

func (provider UserProvider) List() ([]*model.User, error) {
	var users []*model.User
	result := provider.Database.Find(&users)

	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return users, nil
}

func (provider UserProvider) ListPaged(pagination *Pagination[model.User]) (*Pagination[model.User], error) {
	context := paginate(model.User{}, "name", pagination, provider.Database)

	var users []model.User
	result := context.Find(&users)

	if isSqlError(result.Error) {
		return nil, result.Error
	}

	pagination.Rows = users
	return pagination, nil
}

func (provider UserProvider) Get(id uuid.UUID) (*model.User, error) {
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

func (provider UserProvider) GetByName(name string) (*model.User, error) {
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
	user, err := provider.Get(userId)
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

func (provider UserProvider) ListSessions(userId uuid.UUID) ([]*model.Session, error) {
	user, err := provider.Get(userId)
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

func (provider UserProvider) GetBySession(token string) (uuid.UUID, error) {
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

func (provider UserProvider) ValidateByAuthentication(kindName string, configuration string) (*model.User, error) {
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

	return provider.Get(targetAuthentication.UserId)
}

func (provider UserProvider) CreateAuthentication(userId uuid.UUID, kind string, configuration string) (*model.Authentication, error) {
	user, err := provider.Get(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, graphql.NewSafeError("user not found")
	}

	err = json.Unmarshal([]byte(configuration), &struct{}{})
	if err != nil {
		return nil, err
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

func (provider UserProvider) DeleteAuthenticationOfKind(kind string, id uuid.UUID) error {
	var authentication *model.Authentication
	result := provider.Database.Delete(&authentication, "id = ? AND kind = ?", id, kind)
	if isSqlError(result.Error) {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("authentication not found")
	}
	return nil
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
	// Check if migration already ran
	var migration *model.Migration
	result := provider.Database.First(&migration, "name = ?", "default_user")
	if isSqlError(result.Error) {
		return result.Error
	}
	if result.RowsAffected != 0 {
		return nil
	}

	// Create default user
	user, err := provider.Create(
		"admin",
		"",
		"",
		"",
		"default admin user",
		true,
	)
	if err != nil {
		panic(err)
		return err
	}

	// Add all permissions to user
	permissions := make([]string, 0)

	for permission := range Permissions {
		permissions = append(permissions, permission)
	}

	_, err = provider.AddPermission(user.Id, permissions)
	if err != nil {
		return err
	}

	// Create default password
	kind := authenticator.GetKind(authenticator.PASSWORD).(authenticator.Password)

	configuration, err := kind.GenerateConfigurationFromRequestContainer(&authenticator.PasswordRequestContainer{
		Name:     user.Name,
		Password: "admin!",
	})
	if err != nil {
		return err
	}

	configurationJson, err := json.Marshal(configuration)
	if err != nil {
		return err
	}

	_, err = provider.CreateAuthentication(user.Id, authenticator.PASSWORD, string(configurationJson))
	if err != nil {
		return err
	}

	// Create migration
	err = provider.Database.Create(&model.Migration{
		Id:   uuid.New(),
		Name: "default_user",
	}).Error
	if err != nil {
		return err
	}

	return nil
}
