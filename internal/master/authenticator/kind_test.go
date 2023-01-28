package authenticator_test

import (
	"github.com/stretchr/testify/assert"
	"kroseida.org/slixx/internal/master/authenticator"
	"testing"
)

func Test_KindString_PASSWORD(t *testing.T) {
	assert.Equal(t, "password", authenticator.PASSWORD)
}

func Test_GetKind(t *testing.T) {
	kind := authenticator.GetKind(authenticator.PASSWORD)

	if kind == nil {
		t.Error("kind is nil")
	}

	assert.IsType(t, authenticator.Password{}, kind)
}
