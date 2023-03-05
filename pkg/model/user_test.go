package model_test

import (
	"github.com/stretchr/testify/assert"
	"kroseida.org/slixx/pkg/model"
	"testing"
)

func Test_UserGetActivePermissions(t *testing.T) {
	user := model.User{
		Permissions: "[]",
	}

	err := user.AddPermission([]string{"test", "test2"})
	if err != nil {
		t.Error(err)
	}

	isPermittedOr, err := user.IsPermitted([]string{"test", "test2"})
	if err != nil {
		t.Error(err)
	}

	isPermitted, err := user.IsPermitted([]string{"test"})
	if err != nil {
		t.Error(err)
	}

	allPermissions, err := user.GetActivePermissions()
	if err != nil {
		t.Error(err)
	}

	assert.True(t, isPermittedOr)
	assert.True(t, isPermitted)
	assert.Equal(t, 2, len(allPermissions))
}

func Test_UserRemovePermissions(t *testing.T) {
	user := model.User{
		Permissions: "[]",
	}

	err := user.AddPermission([]string{"test", "test2", "flux", "flix", "testa"})
	if err != nil {
		t.Error(err)
	}

	err = user.RemovePermission([]string{"test", "flux"})
	if err != nil {
		t.Error(err)
	}

	allPermissions, err := user.GetActivePermissions()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 3, len(allPermissions))
}
