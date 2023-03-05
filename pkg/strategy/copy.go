package strategy

import (
	"kroseida.org/slixx/pkg/storage"
)

type CopyKind struct {
}

func (kind *CopyKind) Execute(origin storage.Kind, destination []storage.Kind) error {
	files, err := origin.ListFiles()
	if err != nil {
		return err
	}
	// Copy file by file. This should not be a problem for large files based on memory usage.
	for _, file := range files {
		// Right now we are reading just one fill and write one file. Later this could make one file to many files.
		dataMap := make(map[string][]byte)
		dataMap[file], err = origin.Read(file)
		if err != nil {
			return err
		}

		for _, dest := range destination {
			err = dest.Store(dataMap)
			if err != nil {
				return err
			}
		}
		dataMap = nil
	}

	return nil
}

func (kind *CopyKind) Restore(origin storage.Kind, destination []storage.Kind, destinationIndex int) error {
	files, err := destination[destinationIndex].ListFiles()
	if err != nil {
		return err
	}
	dataMap := make(map[string][]byte)
	for _, file := range files {
		dataMap[file], err = origin.Read(file)
		if err != nil {
			return err
		}
	}

	err = origin.Store(dataMap)
	if err != nil {
		return err
	}

	return nil
}
