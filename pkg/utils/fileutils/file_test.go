package fileutils_test

import (
	"github.com/stretchr/testify/assert"
	"kroseida.org/slixx/pkg/utils/fileutils"
	"os"
	"testing"
)

func Test_FileExists(t *testing.T) {
	assert.False(t, fileutils.FileExists("file"))
}

func Test_FileExists_Valid(t *testing.T) {
	os.WriteFile("file", []byte("test"), 0644)
	assert.True(t, fileutils.FileExists("file"))
	os.Remove("file")
}
