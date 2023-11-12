package strategy

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/pkg/storage"
	"reflect"
	"time"
)

type Strategy interface {
	GetName() string
	Initialize(configuration any) error
	// Execute We need the job id to store the backup in the local index
	Execute(jobId uuid.UUID, origin storage.Kind, destination storage.Kind, callback func(BackupStatusUpdate)) (*RawBackupInfo, error)
	Restore(origin storage.Kind, destination storage.Kind, id *uuid.UUID) error
	Parse(configurationJson string) (interface{}, error)
	DefaultConfiguration() interface{}
	ListBackups(destination storage.Kind) ([]*RawBackupInfo, error)
	Close() error
}

type BackupStatusUpdate struct {
	Id         *uuid.UUID `json:"id"`
	JobId      *uuid.UUID `json:"jobId"`
	Percentage float64    `json:"percentage"`
	Message    string     `json:"message"`
	StatusType string     `json:"statusType"`
}

type RawBackupInfo struct {
	Id              *uuid.UUID `json:"id"`
	CreatedAt       time.Time  `json:"createdAt"`
	JobId           *uuid.UUID `json:"jobId"`
	OriginKind      string     `json:"originKind"`
	DestinationKind string     `json:"destinationKind"`
	Strategy        string     `json:"strategy"`
}

var COPY = &CopyStrategy{}

var strategies = map[string]Strategy{
	"COPY": COPY,
}

func ValueOf(name string) Strategy {
	strategy := strategies[name]
	if strategy == nil {
		return nil
	}
	return reflect.New(reflect.TypeOf(strategy).Elem()).Interface().(Strategy)
}

func Values() []Strategy {
	values := make([]Strategy, 0, len(strategies))
	for _, value := range strategies {
		values = append(values, value)
	}
	return values
}
