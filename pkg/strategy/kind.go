package strategy

import "kroseida.org/slixx/pkg/storage"

type Strategy interface {
	Initialize(configuration any) error
	Execute(origin storage.Kind, destination storage.Kind) error
	Restore(origin storage.Kind, destination storage.Kind) error
	Parse(configurationJson string) (interface{}, error)
	DefaultConfiguration() interface{}
}

var strategies = map[string]Strategy{
	"COPY": &CopyStrategy{},
}

func ValueOf(name string) Strategy {
	return strategies[name]
}

func Values() []Strategy {
	values := make([]Strategy, 0, len(strategies))
	for _, value := range strategies {
		values = append(values, value)
	}
	return values
}
