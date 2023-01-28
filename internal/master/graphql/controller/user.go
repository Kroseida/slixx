package controller

import (
	"context"
	"github.com/google/uuid"
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/reactive"
	"kroseida.org/slixx/internal/master/application"
	"kroseida.org/slixx/internal/master/datasource"
	"kroseida.org/slixx/internal/master/datasource/model"
	"kroseida.org/slixx/pkg/dto"
	"time"
)

type User struct {
	Id          uuid.UUID `json:"id",graphql:"id"`
	Name        string    `json:"name",graphql:"name"`
	FirstName   string    `json:"firstName",graphql:"firstName"`
	LastName    string    `json:"lastName",graphql:"lastName"`
	Email       string    `json:"email",graphql:"email"`
	Active      bool      `json:"active",graphql:"active"`
	Permissions string    `json:"permissions",graphql:"permissions"`
	Description string    `json:"description",graphql:"description"`
	CreatedAt   time.Time `json:"createdAt",graphql:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt",graphql:"updatedAt"`
	DeletedAt   time.Time `json:"deletedAt",graphql:"deletedAt"`
}

type Authentication struct {
	Id        uuid.UUID `json:"id",graphql:"id"`
	UserId    uuid.UUID `json:"userId",graphql:"userId"`
	Kind      string    `json:"kind",graphql:"kind"`
	CreatedAt time.Time `json:"createdAt",graphql:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt",graphql:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt",graphql:"deletedAt"`
}

type Session struct {
	Id        uuid.UUID `json:"id",graphql:"id"`
	UserId    uuid.UUID `json:"userId",graphql:"userId"`
	ExpiresAt time.Time `json:"expiresAt",graphql:"expiresAt"`
	CreatedAt time.Time `json:"createdAt",graphql:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt",graphql:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt",graphql:"deletedAt"`
}

type ExposedSession struct {
	Id        uuid.UUID `json:"id",graphql:"id"`
	UserId    uuid.UUID `json:"userId",graphql:"userId"`
	ExpiresAt time.Time `json:"expiresAt",graphql:"expiresAt"`
	CreatedAt time.Time `json:"createdAt",graphql:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt",graphql:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt",graphql:"deletedAt"`
	Token     string    `json:"token",graphql:"token"`
}

func CreateUser(ctx context.Context, args struct {
	Name        string `json:"name",graphql:"name"`
	FirstName   string `json:"firstName",graphql:"firstName"`
	LastName    string `json:"lastName",graphql:"lastName"`
	Email       string `json:"email",graphql:"email"`
	Description string `json:"description",graphql:"description"`
	Active      bool   `json:"active",graphql:"active"`
}) (*User, error) {
	if !IsPermitted(ctx, []string{"user.create"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	user, err := datasource.UserProvider.CreateUser(args.Name, args.Email, args.FirstName, args.LastName, args.Description, args.Active)
	if err != nil {
		application.Logger.Debug(err)
		return nil, err
	}
	var userDto User
	dto.Map(user, &userDto)

	return &userDto, nil
}

func UpdateUserName(ctx context.Context, args struct {
	Id   uuid.UUID `json:"id",graphql:"id"`
	Name string    `json:"name",graphql:"name"`
}) (*User, error) {
	if !IsPermitted(ctx, []string{"user.update"}) && !IsSameUser(ctx, args.Id) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	user, err := datasource.UserProvider.UpdateUserName(args.Id, args.Name)
	if err != nil {
		application.Logger.Debug(err)
		return nil, err
	}
	var userDto User
	dto.Map(&user, &userDto)

	return &userDto, nil
}

func UpdateUserFirstName(ctx context.Context, args struct {
	Id        uuid.UUID `json:"id",graphql:"id"`
	FirstName string    `json:"firstName",graphql:"firstName"`
}) (*User, error) {
	if !IsPermitted(ctx, []string{"user.update"}) && !IsSameUser(ctx, args.Id) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	user, err := datasource.UserProvider.UpdateUserFirstName(args.Id, args.FirstName)
	if err != nil {
		application.Logger.Debug(err)
		return nil, err
	}
	var userDto User
	dto.Map(&user, &userDto)

	return &userDto, nil
}

func UpdateUserLastName(ctx context.Context, args struct {
	Id       uuid.UUID `json:"id",graphql:"id"`
	LastName string    `json:"lastName",graphql:"lastName"`
}) (*User, error) {
	if !IsPermitted(ctx, []string{"user.update"}) && !IsSameUser(ctx, args.Id) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	user, err := datasource.UserProvider.UpdateUserLastName(args.Id, args.LastName)
	if err != nil {
		application.Logger.Debug(err)
		return nil, err
	}
	var userDto User
	dto.Map(&user, &userDto)

	return &userDto, nil
}

func UpdateUserEmail(ctx context.Context, args struct {
	Id    uuid.UUID `json:"id",graphql:"id"`
	Email string    `json:"email",graphql:"email"`
}) (*User, error) {
	if !IsPermitted(ctx, []string{"user.update"}) && !IsSameUser(ctx, args.Id) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	user, err := datasource.UserProvider.UpdateUserEmail(args.Id, args.Email)
	if err != nil {
		application.Logger.Debug(err)
		return nil, err
	}
	var userDto User
	dto.Map(&user, &userDto)

	return &userDto, nil
}

func UpdateUserDescription(ctx context.Context, args struct {
	Id          uuid.UUID `json:"id",graphql:"id"`
	Description string    `json:"description",graphql:"description"`
}) (*User, error) {
	if !IsPermitted(ctx, []string{"user.update"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	user, err := datasource.UserProvider.UpdateUserDescription(args.Id, args.Description)
	if err != nil {
		application.Logger.Debug(err)
		return nil, err
	}
	var userDto User
	dto.Map(&user, &userDto)

	return &userDto, nil
}

func UpdateUserActive(ctx context.Context, args struct {
	Id     uuid.UUID `json:"id",graphql:"id"`
	Active bool      `json:"active",graphql:"active"`
}) (*User, error) {
	if !IsPermitted(ctx, []string{"user.update"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	user, err := datasource.UserProvider.UpdateUserActive(args.Id, args.Active)
	if err != nil {
		application.Logger.Debug(err)
		return nil, err
	}
	var userDto User
	dto.Map(&user, &userDto)

	return &userDto, nil
}

func AddUserPermission(ctx context.Context, args struct {
	Id          uuid.UUID `json:"id",graphql:"id"`
	Permissions []string  `json:"permissions",graphql:"permissions"`
}) (*User, error) {
	if !IsPermitted(ctx, []string{"user.permission.update"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	user, err := datasource.UserProvider.AddUserPermission(args.Id, args.Permissions)
	if err != nil {
		application.Logger.Debug(err)
		return nil, err
	}

	var userDto User
	dto.Map(&user, &userDto)

	return &userDto, nil
}

func RemoveUserPermission(ctx context.Context, args struct {
	Id          uuid.UUID `json:"id",graphql:"id"`
	Permissions []string  `json:"permissions",graphql:"permissions"`
}) (*User, error) {
	if !IsPermitted(ctx, []string{"user.permission.update"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	user, err := datasource.UserProvider.RemoveUserPermission(args.Id, args.Permissions)
	if err != nil {
		application.Logger.Debug(err)
		return nil, err
	}

	var userDto User
	dto.Map(&user, &userDto)

	return &userDto, nil
}

func CreatePasswordAuthentication(ctx context.Context, args struct {
	Id       uuid.UUID `json:"id",graphql:"id"`
	Password string    `json:"password",graphql:"password"`
}) (*Authentication, error) {
	if !IsPermitted(ctx, []string{"user.update"}) && !IsSameUser(ctx, args.Id) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	authentication, err := datasource.UserProvider.CreatePasswordAuthentication(args.Id, args.Password)
	if err != nil {
		application.Logger.Debug(err)
		return nil, err
	}
	var authenticationDto Authentication
	dto.Map(authentication, &authenticationDto)

	return &authenticationDto, nil
}

func Authenticate(ctx context.Context, args struct {
	Name     string `json:"name",graphql:"name"`
	Password string `json:"password",graphql:"password"`
}) (*ExposedSession, error) {
	reactive.InvalidateAfter(ctx, 5*time.Second)
	session, err := datasource.UserProvider.AuthenticatePassword(args.Name, args.Password)
	if err != nil {
		application.Logger.Debug(err)
		return nil, err
	}
	var sessionDto ExposedSession
	dto.Map(session, &sessionDto)

	return &sessionDto, nil
}

func GetUsers(ctx context.Context) ([]*User, error) {
	if !IsPermitted(ctx, []string{"user.view"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	users, err := datasource.UserProvider.GetUsers()
	if err != nil {
		application.Logger.Debug(err)
		return nil, err
	}
	var userDtos []*User
	dto.Map(&users, &userDtos)

	return userDtos, nil
}

func GetUser(ctx context.Context, args struct {
	Id uuid.UUID `json:"id",graphql:"id"`
}) (*User, error) {
	if !IsPermitted(ctx, []string{"user.view"}) && !IsSameUser(ctx, args.Id) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	user, err := datasource.UserProvider.GetUser(args.Id)
	if err != nil {
		application.Logger.Debug(err)
		return nil, err
	}
	var userDto User
	dto.Map(&user, &userDto)

	return &userDto, nil
}

func GetLocalUser(ctx context.Context, _ struct{}) (*User, error) {
	reactive.InvalidateAfter(ctx, 5*time.Second)
	if ctx.Value("user").(*model.User) == nil {
		return nil, nil
	}
	user, err := datasource.UserProvider.GetUser(ctx.Value("user").(*model.User).Id)
	if err != nil {
		application.Logger.Debug(err)
		return nil, err
	}
	var userDto User
	dto.Map(&user, &userDto)

	return &userDto, nil
}
