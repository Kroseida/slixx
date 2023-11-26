package strategy

import (
	"encoding/json"
	"github.com/google/uuid"
	"kroseida.org/slixx/pkg/storage"
	"kroseida.org/slixx/pkg/strategy/statustype"
	"kroseida.org/slixx/pkg/utils/fileutils"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var SlixxDirectory = ".slixx"
var BackupInfoDirectory = SlixxDirectory + "/backups/"

type CopyStrategy struct {
	Configuration *CopyStrategyConfiguration
}

type CopyStrategyConfiguration struct {
	BlockSize int `json:"blockSize" slixx:"LONG" default:"1073741824"` // 1GB in bytes
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

func (strategy *CopyStrategy) Execute(jobId uuid.UUID, origin storage.Kind, destination storage.Kind, callback func(BackupStatusUpdate)) (*RawBackupInfo, error) {
	rawFiles, id, err := strategy.handleInitialize(origin, destination, callback)
	if err != nil {
		return nil, err
	}
	var parallels [][]fileutils.FileInfo
	var parallelStatus = make([]float64, strategy.Configuration.Parallel)
	var parallelError = make(chan error)
	var parallelFinished = make([]bool, strategy.Configuration.Parallel)
	var parallelProceededBytes = make([]uint64, strategy.Configuration.Parallel)

	parallels, sizes := fileutils.SplitArrayBySize(rawFiles, strategy.Configuration.Parallel)

	callback(BackupStatusUpdate{
		Id:         &id,
		Percentage: 0,
		Message:    "Opening " + strconv.Itoa(len(parallels)) + " parallels",
		StatusType: statustype.Info,
	})

	for index, _ := range parallels {
		atomicIndex := index
		parallelFinished[atomicIndex] = false
		parallelStatus[atomicIndex] = 0
		go func() {
			files := parallels[atomicIndex]

			originCopy := reflect.New(reflect.TypeOf(origin).Elem()).Interface().(storage.Kind)
			err = originCopy.Initialize(origin.GetConfiguration())
			if err != nil {
				parallelError <- err
				return
			}

			destinationCopy := reflect.New(reflect.TypeOf(destination).Elem()).Interface().(storage.Kind)
			err := destinationCopy.Initialize(destination.GetConfiguration())
			if err != nil {
				parallelError <- err
				return
			}

			for _, file := range files {
				if file.Directory {
					continue
				}
				tryCopy(strategy, originCopy, destinationCopy, file, id, parallelError)
				parallelProceededBytes[atomicIndex] += file.Size
				if sizes[atomicIndex] == 0 || parallelProceededBytes[atomicIndex] == 0 {
					parallelStatus[atomicIndex] = 0
				} else {
					parallelStatus[atomicIndex] = float64(parallelProceededBytes[atomicIndex]) / float64(sizes[atomicIndex])
				}
			}

			// Close the connections
			originCopy.Close()
			destinationCopy.Close()
			parallelFinished[atomicIndex] = true
		}()
	}
	err = strategy.handleBackupWatchdog(parallelStatus, parallelFinished, callback, id, parallelError)
	if err != nil {
		return nil, err
	}

	return strategy.handleIndexingBackup(id, jobId, origin, destination, callback)
}

func tryCopy(strategy *CopyStrategy, originCopy storage.Kind, destinationCopy storage.Kind, file fileutils.FileInfo,
	id uuid.UUID, parallelError chan error) error {
	retries := 0
	for {
		copyErr := strategy.copy(originCopy, destinationCopy, file, id.String())
		if copyErr == nil {
			break
		}
		retries++
		if retries > 15 {
			parallelError <- copyErr
			return copyErr
		}
		// Wait 3 second before retrying so that we don't overload the server
		time.Sleep(3 * time.Second)
	}
	return nil
}

func (strategy *CopyStrategy) handleInitialize(origin storage.Kind, destination storage.Kind, callback func(BackupStatusUpdate)) ([]fileutils.FileInfo, uuid.UUID, error) {
	id := uuid.New()
	callback(BackupStatusUpdate{
		Id:         &id,
		Percentage: 0,
		Message:    "Initializing strategy",
		StatusType: statustype.Info,
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
				return nil, id, err
			}
		}
	}

	callback(BackupStatusUpdate{
		Id:         &id,
		Percentage: 0,
		Message:    "Detected " + strconv.Itoa(len(rawFiles)) + " files",
		StatusType: statustype.Info,
	})

	if err != nil {
		callback(BackupStatusUpdate{
			Id:         &id,
			Percentage: 0,
			Message:    err.Error(),
			StatusType: statustype.Error,
		})
		return nil, id, err
	}

	return rawFiles, id, nil
}

func (strategy *CopyStrategy) handleBackupWatchdog(parallelStatus []float64, parallelFinished []bool,
	callback func(BackupStatusUpdate), id uuid.UUID, parallelError chan error) error {
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
			Message:    "Copying files from origin to destination",
			StatusType: statustype.Info,
		})

		var err error
		// Check for errors in the parallelError channel
		select {
		case res := <-parallelError:
			err = res
		case <-time.After(1000 * time.Millisecond):
			err = nil
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
	}
	return nil
}

func (strategy *CopyStrategy) handleIndexingBackup(id uuid.UUID, jobId uuid.UUID, origin storage.Kind, destination storage.Kind,
	callback func(BackupStatusUpdate)) (*RawBackupInfo, error) {

	callback(BackupStatusUpdate{
		Id:         &id,
		Percentage: 100,
		Message:    "Creating backup info file on destination",
		StatusType: statustype.Info,
	})

	handleCreateSlixxDirectories(destination)

	// This is the raw backup info that we will store in the destination storage, so that we can restore it later e.g
	// after supervisor get corrupted data or something like that we can restore this backup info and continue
	rawBackupInfo := RawBackupInfo{
		Id:              &id,
		CreatedAt:       time.Now(),
		JobId:           &jobId,
		OriginKind:      origin.GetName(),      // Store origin kind so that we can restore it later
		DestinationKind: destination.GetName(), // Store destination kind so that we can restore it later
		Strategy:        strategy.GetName(),    // Store strategy so that we can restore it later
	}
	rawBackupInfoBytes, err := json.Marshal(rawBackupInfo)
	if err != nil {
		callback(BackupStatusUpdate{
			Id:         &id,
			Percentage: 0,
			Message:    err.Error(),
			StatusType: statustype.Error,
		})
		return nil, err
	}

	// Store backup info file in destination so that we have some information about the backup
	err = destination.Store(
		fileutils.FixedPathName(BackupInfoDirectory+"/"+id.String()),
		rawBackupInfoBytes,
		0,
	)
	if err != nil {
		callback(BackupStatusUpdate{
			Id:         &id,
			Percentage: 0,
			Message:    err.Error(),
			StatusType: statustype.Error,
		})
		return nil, err
	}

	callback(BackupStatusUpdate{
		Id:         &id,
		Percentage: 100,
		Message:    "FINISHED",
		StatusType: statustype.Finished,
	})

	return &rawBackupInfo, nil
}

func handleCreateSlixxDirectories(destination storage.Kind) {
	// Create .slixx directory and ignore errors
	destination.CreateDirectory(fileutils.FixedPathName(SlixxDirectory))
	destination.CreateDirectory(fileutils.FixedPathName(BackupInfoDirectory))
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
	handleCreateSlixxDirectories(destination)

	files, err := destination.ListFiles(BackupInfoDirectory)
	if err != nil {
		return nil, err
	}

	backupList := make([]*RawBackupInfo, 0, len(files))
	for _, file := range files {
		parsedId, err := uuid.Parse(strings.TrimPrefix(file.RelativePath, BackupInfoDirectory))
		if err != nil {
			return nil, err
		}
		backupList = append(backupList, &RawBackupInfo{
			Id:        &parsedId,
			CreatedAt: time.Unix(file.CreatedAt, 0),
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
