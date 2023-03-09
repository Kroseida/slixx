package strategy

import "kroseida.org/slixx/pkg/storage"

type Kind interface {
	Execute(origin storage.Kind, destination []storage.Kind) error
	Restore(origin storage.Kind, destination []storage.Kind, destinationIndex int) error
}

var COPY Kind = &CopyKind{}

func ValueOf(name string) Kind {
	if name == "COPY" {
		return COPY
	}
	return nil
}

func Values() []Kind {
	return []Kind{
		COPY,
	}
}
