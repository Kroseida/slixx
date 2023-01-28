package utils_test

import (
	"github.com/stretchr/testify/assert"
	"kroseida.org/slixx/pkg/utils"
	"os"
	"testing"
)

func Test_FileExists(t *testing.T) {
	assert.False(t, utils.FileExists("file"))
}

func Test_FileExists_Valid(t *testing.T) {
	os.WriteFile("file", []byte("test"), 0644)
	assert.True(t, utils.FileExists("file"))
	os.Remove("file")
}
