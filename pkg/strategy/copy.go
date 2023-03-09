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
		err = kind.copy(origin, destination, file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (kind *CopyKind) Restore(origin storage.Kind, destination []storage.Kind, destinationIndex int) error {
	files, err := destination[destinationIndex].ListFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		err = kind.copy(destination[destinationIndex], []storage.Kind{origin}, file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (kind *CopyKind) copy(origin storage.Kind, destination []storage.Kind, file string) error {
	blockSize := 1024 // TODO: Make this configurable

	// Read File Size
	size, err := origin.Size(file)
	if err != nil {
		return err
	}

	iterations := int(size) / blockSize
	lastBlockSize := int(size) % blockSize
	if lastBlockSize != 0 {
		iterations++
	}

	for i := 0; i < iterations; i++ {
		readSize := blockSize
		if i == iterations-1 && lastBlockSize != 0 {
			readSize = lastBlockSize
		}
		data, err := origin.Read(file, uint64(i*blockSize), uint64(readSize))
		if err != nil {
			return err
		}
		for _, dest := range destination {
			err = dest.Store(file, data, uint64(i*blockSize))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
