package strategy

import (
	"encoding/json"
	"kroseida.org/slixx/pkg/storage"
)

type CopyStrategy struct {
	Configuration *CopyStrategyConfiguration
}

type CopyStrategyConfiguration struct {
	blockSize int
}

func (strategy *CopyStrategy) Initialize(rawConfiguration any) error {
	configuration := rawConfiguration.(*CopyStrategyConfiguration)

	strategy.Configuration = configuration

	return nil
}

func (strategy *CopyStrategy) Execute(origin storage.Kind, destination storage.Kind) error {
	files, err := origin.ListFiles()
	if err != nil {
		return err
	}
	// Copy file by file. This should not be a problem for large files based on memory usage.
	for _, file := range files {
		err = strategy.copy(origin, destination, file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (strategy *CopyStrategy) Restore(origin storage.Kind, destination storage.Kind) error {
	files, err := destination.ListFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		err = strategy.copy(destination, origin, file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (strategy *CopyStrategy) copy(origin storage.Kind, destination storage.Kind, file string) error {
	// Read File Size
	size, err := origin.Size(file)
	if err != nil {
		return err
	}

	iterations := int(size) / strategy.Configuration.blockSize
	lastBlockSize := int(size) % strategy.Configuration.blockSize
	if lastBlockSize != 0 {
		iterations++
	}

	for index := 0; index < iterations; index++ {
		readSize := strategy.Configuration.blockSize
		if index == iterations-1 && lastBlockSize != 0 {
			readSize = lastBlockSize
		}

		data, err := origin.Read(file, uint64(index*strategy.Configuration.blockSize), uint64(readSize))
		if err != nil {
			return err
		}

		err = destination.Store(file, data, uint64(index*strategy.Configuration.blockSize))
		if err != nil {
			return err
		}
	}
	return nil
}

func (strategy *CopyStrategy) Parse(configurationJson string) (interface{}, error) {
	var configuration CopyStrategyConfiguration
	err := json.Unmarshal([]byte(configurationJson), &configuration)
	if err != nil {
		return nil, err
	}
	return &configuration, nil
}

func (strategy *CopyStrategy) DefaultConfiguration() interface{} {
	return &CopyStrategyConfiguration{}
}
