package schedule

import (
	"reflect"
)

type Kind interface {
	GetName() string
	Initialize(configuration any, callback func()) error
	Deactivate() error
	Parse(configurationJson string) (interface{}, error)
	DefaultConfiguration() interface{}
}

var kinds = map[string]Kind{
	"CRON": &CronKind{},
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
