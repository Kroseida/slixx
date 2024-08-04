package strategy

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/pkg/storage"
	"kroseida.org/slixx/pkg/utils/parallel"
	"reflect"
	"time"
)

type Strategy interface {
	// GetName Get the name of the strategy (used to identify it)
	GetName() string
	// Initialize Initialize the strategy with the configuration, this is called before any other method
	Initialize(configuration any) error
	// Execute The main method of the strategy execute a backup from the origin to the destination storage (this is called when a backup is requested)
	Execute(job *parallel.RunningJob, origin storage.Kind, destination storage.Kind) (*RawBackupInfo, error)
	// Restore Restore a backup from the destination to the origin storage (this is called when a restore is requested)
	Restore(job *parallel.RunningJob, origin storage.Kind, destination storage.Kind, id *uuid.UUID) error
	// Delete Delete a backup from the destination storage (this is called when a delete is requested)
	Delete(job *parallel.RunningJob, destination storage.Kind, id *uuid.UUID) error
	// Parse Parse the configuration of the strategy from a json string to a struct
	Parse(configurationJson string) (interface{}, error)
	// DefaultConfiguration Get the DefaultConfiguration Get the default configuration of the strategy
	DefaultConfiguration() interface{}
	// ListBackups List the backups stored in the destination storage, this is not always the source of truth
	// We can later on make it so that the supervisor is the source of truth and the satellite just sends the backups infos
	ListBackups(destination storage.Kind) ([]*RawBackupInfo, error)
	// Close Close the strategy and all its resources
	Close() error
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
var INCREMENTAL = &IncrementalStrategy{}

var strategies = map[string]Strategy{
	"COPY":        COPY,
	"INCREMENTAL": INCREMENTAL,
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
