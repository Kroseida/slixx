package storage

import (
	"kroseida.org/slixx/pkg/utils/fileutils"
	"reflect"
)

type Kind interface {
	GetName() string
	Initialize(configuration any) error
	Store(name string, data []byte, offset uint64) error
	FileInfo(name string) (fileutils.FileInfo, error)
	CreateDirectory(name string) error
	ListFiles(directory string) ([]fileutils.FileInfo, error)
	Size(file string) (uint64, error)
	Read(file string, offset uint64, size uint64) ([]byte, error)
	Delete(file string) error
	DeleteDirectory(directory string) error
	Parse(configurationJson string) (interface{}, error)
	DefaultConfiguration() interface{}
	CanStore() bool
	CanRead() bool
	Close() error
	GetConfiguration() any
}

var kinds = map[string]Kind{
	"FTP":  &FtpKind{},
	"SFTP": &SFtpKind{},
}

func ValueOf(name string) Kind {
	kind := kinds[name]
	if kind == nil {
		return nil
	}
	return reflect.New(reflect.TypeOf(kind).Elem()).Interface().(Kind)
}

func Values() []Kind {
	values := make([]Kind, 0, len(kinds))
	for _, value := range kinds {
		values = append(values, value)
	}
	return values
}
