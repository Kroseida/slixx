package strategy

import (
	"encoding/json"
	"github.com/google/uuid"
	"kroseida.org/slixx/pkg/storage"
	"kroseida.org/slixx/pkg/strategy/statustype"
	"kroseida.org/slixx/pkg/utils/fileutils"
	"reflect"
	"strings"
	"time"
)

type CopyStrategy struct {
	Configuration *CopyStrategyConfiguration
}

type CopyStrategyConfiguration struct {
	BlockSize int `json:"blockSize" slixx:"LONG" default:"4096"`
	Parallel  int `json:"parallel" slixx:"LONG" default:"4"`
}

func (strategy *CopyStrategy) GetName() string {
	return "COPY"
}

func (strategy *CopyStrategy) Initialize(rawConfiguration any) error {
	configuration := rawConfiguration.(*CopyStrategyConfiguration)

	strategy.Configuration = configuration

	return nil
}

func (strategy *CopyStrategy) Execute(origin storage.Kind, destination storage.Kind, callback func(BackupStatusUpdate)) error {
	id := uuid.New()
	callback(BackupStatusUpdate{
		Id:         &id,
		Percentage: 0,
		Message:    "Initializing",
		StatusType: "INFO",
	})
	rawFiles, err := origin.ListFiles("")
	for _, file := range rawFiles {
		if file.Directory {
			err := destination.CreateDirectory(id.String() + file.RelativePath)
			if err != nil {
				callback(BackupStatusUpdate{
					Id:         &id,
					Percentage: 0,
					Message:    err.Error(),
					StatusType: statustype.Error,
				})
				return err
			}
		}
	}

	if err != nil {
		callback(BackupStatusUpdate{
			Id:         &id,
			Percentage: 0,
			Message:    err.Error(),
			StatusType: statustype.Error,
		})
		return err
	}

	var parallels [][]fileutils.FileInfo
	var parallelStatus = make([]float64, strategy.Configuration.Parallel)
	var parallelError = make([]error, strategy.Configuration.Parallel)
	var parallelFinished = make([]bool, strategy.Configuration.Parallel)

	parallels = fileutils.SplitArrayBySize(rawFiles, strategy.Configuration.Parallel)

	for index, _ := range parallels {
		atomicIndex := index
		parallelFinished[atomicIndex] = false
		parallelStatus[atomicIndex] = 0
		go func() {
			files := parallels[atomicIndex]

			originCopy := reflect.New(reflect.TypeOf(origin).Elem()).Interface().(storage.Kind)
			err = originCopy.Initialize(destination.GetConfiguration())
			if err != nil {
				parallelError[atomicIndex] = err
				return
			}

			destinationCopy := reflect.New(reflect.TypeOf(destination).Elem()).Interface().(storage.Kind)
			err := destinationCopy.Initialize(destination.GetConfiguration())
			if err != nil {
				parallelError[atomicIndex] = err
				return
			}

			for index, file := range files {
				if file.Directory {
					continue
				}
				retries := 0
				for {
					copyErr := strategy.copy(originCopy, destinationCopy, file, id.String())
					if copyErr == nil {
						break
					}
					retries++
					if retries > 15 {
						parallelError[atomicIndex] = err
						return
					}
					// Wait 1 second before retrying so that we don't overload the server
					time.Sleep(1 * time.Second)
				}
				parallelStatus[atomicIndex] = float64(index) / float64(len(files))
			}

			// Close the connections
			originCopy.Close()
			destinationCopy.Close()
			parallelFinished[atomicIndex] = true
		}()
	}

	for {
		allFinished := true
		for _, finished := range parallelFinished {
			if !finished {
				allFinished = false
				break
			}
		}
		if allFinished {
			break
		}

		percentage := 0.0
		for _, status := range parallelStatus {
			percentage += status
		}
		percentage /= float64(len(parallelStatus))
		callback(BackupStatusUpdate{
			Id:         &id,
			Percentage: percentage * 100,
			Message:    "COPYING",
			StatusType: statustype.Info,
		})

		for _, err := range parallelError {
			if err != nil {
				callback(BackupStatusUpdate{
					Id:         &id,
					Percentage: 0,
					Message:    err.Error(),
					StatusType: statustype.Error,
				})
				return err
			}
		}

		time.Sleep(5 * time.Second)
	}

	callback(BackupStatusUpdate{
		Id:         &id,
		Percentage: 100,
		Message:    "INDEXING",
		StatusType: statustype.Info,
	})

	// Store backup info file in destination so that we have some information about the backup
	err = destination.Store(".slixx/backups/"+id.String(), []byte{1}, 0)
	if err != nil {
		callback(BackupStatusUpdate{
			Id:         &id,
			Percentage: 0,
			Message:    err.Error(),
			StatusType: statustype.Error,
		})
		return err
	}

	callback(BackupStatusUpdate{
		Id:         &id,
		Percentage: 100,
		Message:    "FINISHED",
		StatusType: statustype.Finished,
	})
	return nil
}

func (strategy *CopyStrategy) Restore(origin storage.Kind, destination storage.Kind, id *uuid.UUID) error {
	files, err := destination.ListFiles(id.String())
	if err != nil {
		return err
	}
	for _, file := range files {
		if !file.Directory {
			err = strategy.copy(origin, destination, file, id.String())
		} else {
			err = destination.CreateDirectory(id.String() + "/" + file.FullDirectory)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func (strategy *CopyStrategy) ListBackups(destination storage.Kind) ([]*RawBackupInfo, error) {
	files, err := destination.ListFiles(".slixx/info")
	if err != nil {
		return nil, err
	}

	backupList := make([]*RawBackupInfo, 0, len(files))
	for _, file := range files {
		parsedId, err := uuid.Parse(strings.TrimPrefix(file.RelativePath, ".slixx/info/"))
		if err != nil {
			return nil, err
		}
		backupList = append(backupList, &RawBackupInfo{
			Id:        &parsedId,
			CreatedAt: file.CreatedAt,
		})
	}

	return backupList, nil
}

func (strategy *CopyStrategy) copy(origin storage.Kind, destination storage.Kind, file fileutils.FileInfo, storePrefix string) error {
	if file.Directory {
		return nil
	}

	// Read File Size
	size, err := origin.Size(file.FullDirectory)
	if err != nil {
		return err
	}

	iterations := int(size) / strategy.Configuration.BlockSize
	lastBlockSize := int(size) % strategy.Configuration.BlockSize

	if lastBlockSize != 0 {
		iterations++
	}
	if iterations == 0 {
		err = destination.Store(storePrefix+"/"+file.RelativePath, []byte{}, uint64(0))
		if err != nil {
			return err
		}
	}

	for index := 0; index < iterations; index++ {
		readSize := strategy.Configuration.BlockSize
		if index == iterations-1 && lastBlockSize != 0 {
			readSize = lastBlockSize
		}

		data, err := origin.Read(file.FullDirectory, uint64(index*strategy.Configuration.BlockSize), uint64(readSize))
		if err != nil {
			return err
		}

		err = destination.Store(storePrefix+"/"+file.RelativePath, data, uint64(index*strategy.Configuration.BlockSize))
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
	return reflect.New(reflect.TypeOf(CopyStrategyConfiguration{})).Interface()
}

func (strategy *CopyStrategy) Close() error {
	return nil
}
