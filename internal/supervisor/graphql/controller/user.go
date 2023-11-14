package controller

import (
	"context"
	"github.com/google/uuid"
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/reactive"
	"kroseida.org/slixx/internal/supervisor/datasource/provider"
	userService "kroseida.org/slixx/internal/supervisor/service/user"
	"kroseida.org/slixx/pkg/dto"
	"kroseida.org/slixx/pkg/model"
	"time"
)

type User struct {
	Id          uuid.UUID `json:"id" graphql:"id"`
	Name        string    `json:"name" graphql:"name"`
	FirstName   string    `json:"firstName" graphql:"firstName"`
	LastName    string    `json:"lastName" graphql:"lastName"`
	Email       string    `json:"email" graphql:"email"`
	Active      bool      `json:"active" graphql:"active"`
	Permissions string    `json:"permissions" graphql:"permissions"`
	Description string    `json:"description" graphql:"description"`
	CreatedAt   time.Time `json:"createdAt" graphql:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" graphql:"updatedAt"`
	DeletedAt   time.Time `json:"deletedAt" graphql:"deletedAt"`
}

type UsersPage struct {
	Rows []User `json:"rows" graphql:"rows"`
	Page
}

type Authentication struct {
	Id        uuid.UUID `json:"id" graphql:"id"`
	UserId    uuid.UUID `json:"userId" graphql:"userId"`
	Kind      string    `json:"kind" graphql:"kind"`
	CreatedAt time.Time `json:"createdAt" graphql:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" graphql:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt" graphql:"deletedAt"`
}

type Session struct {
	Id        uuid.UUID `json:"id" graphql:"id"`
	UserId    uuid.UUID `json:"userId" graphql:"userId"`
	ExpiresAt time.Time `json:"expiresAt" graphql:"expiresAt"`
	CreatedAt time.Time `json:"createdAt" graphql:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" graphql:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt" graphql:"deletedAt"`
}

type ExposedSession struct {
	Id        uuid.UUID `json:"id" graphql:"id"`
	UserId    uuid.UUID `json:"userId" graphql:"userId"`
	ExpiresAt time.Time `json:"expiresAt" graphql:"expiresAt"`
	CreatedAt time.Time `json:"createdAt" graphql:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" graphql:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt" graphql:"deletedAt"`
	Token     string    `json:"token" graphql:"token"`
}

type CreateUserDto struct {
	Name        string `json:"name" graphql:"name"`
	FirstName   string `json:"firstName" graphql:"firstName"`
	LastName    string `json:"lastName" graphql:"lastName"`
	Email       string `json:"email" graphql:"email"`
	Description string `json:"description" graphql:"description"`
	Active      bool   `json:"active" graphql:"active"`
}

func CreateUser(ctx context.Context, args CreateUserDto) (*User, error) {
	if !IsPermitted(ctx, []string{"user.create"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	user, err := userService.Create(args.Name, args.Email, args.FirstName, args.LastName, args.Description, args.Active)
	if err != nil {
		return nil, err
	}
	var userDto User
	dto.Map(user, &userDto)

	return &userDto, nil
}

type UpdateUserDto struct {
	Id          uuid.UUID `json:"id" graphql:"id"`
	Name        *string   `json:"name" graphql:"name"`
	FirstName   *string   `json:"firstName" graphql:"firstName"`
	LastName    *string   `json:"lastName" graphql:"lastName"`
	Active      *bool     `json:"active" graphql:"active"`
	Description *string   `json:"description" graphql:"description"`
	Email       *string   `json:"email" graphql:"email"`
}

func UpdateUser(ctx context.Context, args UpdateUserDto) (*User, error) {
	if !IsPermitted(ctx, []string{"user.update"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	user, err := userService.Update(args.Id, args.Name, args.FirstName, args.LastName, args.Active, args.Description, args.Email)
	if err != nil {
		return nil, err
	}
	var userDto User
	dto.Map(&user, &userDto)

	return &userDto, nil
}

type AddUserPermissionDto struct {
	Id          uuid.UUID `json:"id" graphql:"id"`
	Permissions []string  `json:"permissions" graphql:"permissions"`
}

func AddUserPermission(ctx context.Context, args AddUserPermissionDto) (*User, error) {
	if !IsPermitted(ctx, []string{"user.permission.update"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	user, err := userService.AddPermission(args.Id, args.Permissions)
	if err != nil {
		return nil, err
	}

	var userDto User
	dto.Map(&user, &userDto)

	return &userDto, nil
}

type RemoveUserPermissionDto struct {
	Id          uuid.UUID `json:"id" graphql:"id"`
	Permissions []string  `json:"permissions" graphql:"permissions"`
}

func RemoveUserPermission(ctx context.Context, args RemoveUserPermissionDto) (*User, error) {
	if !IsPermitted(ctx, []string{"user.permission.update"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	user, err := userService.RemovePermission(args.Id, args.Permissions)
	if err != nil {
		return nil, err
	}

	var userDto User
	dto.Map(&user, &userDto)

	return &userDto, nil
}

type UpdateUserPasswordDto struct {
	Id       uuid.UUID `json:"id" graphql:"id"`
	Password string    `json:"password" graphql:"password"`
}

func CreatePasswordAuthentication(ctx context.Context, args UpdateUserPasswordDto) (*Authentication, error) {
	if !IsPermitted(ctx, []string{"user.update"}) && !IsSameUser(ctx, args.Id) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	authentication, err := userService.CreatePasswordAuthentication(args.Id, args.Password)
	if err != nil {
		return nil, err
	}
	var authenticationDto Authentication
	dto.Map(authentication, &authenticationDto)

	return &authenticationDto, nil
}

type PasswordAuthenticationDto struct {
	Name     string `json:"name" graphql:"name"`
	Password string `json:"password" graphql:"password"`
}

func Authenticate(ctx context.Context, args PasswordAuthenticationDto) (*ExposedSession, error) {
	reactive.InvalidateAfter(ctx, 5*time.Second)
	session, err := userService.AuthenticatePassword(args.Name, args.Password)
	if err != nil {
		return nil, err
	}
	var sessionDto ExposedSession
	dto.Map(session, &sessionDto)

	return &sessionDto, nil
}

func GetUsers(ctx context.Context, args PageArgs) (*UsersPage, error) {
	if !IsPermitted(ctx, []string{"user.view"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)

	var pagination provider.Pagination[model.User]
	dto.Map(&args, &pagination)

	users, err := userService.GetPaged(&pagination)
	if err != nil {
		return nil, err
	}

	var userDtos UsersPage
	dto.Map(&users, &userDtos)

	return &userDtos, nil
}

type GetUserDto struct {
	Id uuid.UUID `json:"id" graphql:"id"`
}

func GetUser(ctx context.Context, args GetUserDto) (*User, error) {
	if !IsPermitted(ctx, []string{"user.view"}) && !IsSameUser(ctx, args.Id) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	user, err := userService.Get(args.Id)
	if err != nil {
		return nil, err
	}
	var userDto *User
	dto.Map(&user, &userDto)

	return userDto, nil
}

type DeleteUserDto struct {
	Id uuid.UUID `json:"id" graphql:"id"`
}

func DeleteUser(ctx context.Context, args DeleteUserDto) (*User, error) {
	if !IsPermitted(ctx, []string{"user.delete"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	user, err := userService.Delete(args.Id)
	if err != nil {
		return nil, err
	}
	var userDto User
	dto.Map(&user, &userDto)

	return &userDto, nil
}

func GetLocalUser(ctx context.Context) (*User, error) {
	reactive.InvalidateAfter(ctx, 5*time.Second)
	if ctx.Value("user").(*model.User) == nil {
		return nil, nil
	}
	user, err := userService.Get(ctx.Value("user").(*model.User).Id)
	if err != nil {
		return nil, err
	}
	var userDto User
	dto.Map(&user, &userDto)

	return &userDto, nil
}

type PermissionDto struct {
	Value string `json:"value" graphql:"value"`
	Name  string `json:"name" graphql:"name"`
}

func GetPermissions() ([]PermissionDto, error) {
	permissionDtos := make([]PermissionDto, 0)
	for value, name := range provider.Permissions {
		permissionDtos = append(permissionDtos, PermissionDto{
			Value: value,
			Name:  name,
		})
	}
	return permissionDtos, nil
}
