package controller

import (
	"context"
	"github.com/google/uuid"
	"kroseida.org/slixx/pkg/model"
)

func IsPermitted(ctx context.Context, permissions []string) bool {
	if ctx.Value("user") == nil {
		return false
	}
	user := ctx.Value("user").(*model.User)
	if user == nil {
		return false
	}
	permitted, err := user.IsPermitted(permissions)
	if err != nil {
		return false
	}
	return permitted
}

func IsSameUser(ctx context.Context, id uuid.UUID) bool {
	if ctx.Value("user") == nil {
		return false
	}
	user := ctx.Value("user").(*model.User)
	if user == nil {
		return false
	}
	return user.Id == id
}