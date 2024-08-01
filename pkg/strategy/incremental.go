package strategy

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/gorm/utils"
	"kroseida.org/slixx/pkg/statustype"
	"kroseida.org/slixx/pkg/storage"
	"kroseida.org/slixx/pkg/utils/fileutils"
	"kroseida.org/slixx/pkg/utils/parallel"
	"reflect"
	"strconv"
	"time"
)

const BlocksDirectory = "blocks/"

type IncrementalStrategy struct {
	Configuration *IncrementalStrategyConfiguration
}

type IncrementalStrategyConfiguration struct {
	BlockSize int64 `json:"blockSize" slixx:"LONG" default:"104857600"` // 100MB
	Parallel  int   `json:"parallel" slixx:"LONG" default:"4"`
}

func (_ *IncrementalStrategy) Hash(data []byte) string {
	hasher := sha512.New()
	hasher.Write(data)
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func (strategy *IncrementalStrategy) GetName() string {
	return "INCREMENTAL"
}

func (strategy *IncrementalStrategy) Initialize(rawConfiguration any) error {
	configuration := rawConfiguration.(*IncrementalStrategyConfiguration)

	strategy.Configuration = configuration

	return nil
}

func (strategy *IncrementalStrategy) Execute(jobId uuid.UUID, origin storage.Kind, destination storage.Kind, callback func(StatusUpdate)) (*RawBackupInfo, error) {
	handleCreateSlixxDirectories(destination)

	destination.CreateDirectory(fileutils.FixedPathName(BlocksDirectory))

	backupId := uuid.New()
	rawFiles, err := strategy.handleInitialize(backupId.String(), "", origin, destination, callback)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	blockFiles, err := destination.ListFiles(BlocksDirectory)
	if err != nil {
		return nil, err
	}

	var blocks = make([]string, len(blockFiles))
	for index, blockFile := range blockFiles {
		blockHash := blockFile.RelativePath
		if blockHash[0] == '/' || blockHash[0] == '\\' {
			blockHash = blockHash[1:]
		}
		blocks[index] = blockHash
	}

	parallelFiles, sizes := fileutils.SplitArrayBySize(rawFiles, strategy.Configuration.Parallel)

	parallelExecutor := parallel.NewExecutor(
		parallelFiles,
		func(executorStatus parallel.ExecutorStatus) {
			callback(StatusUpdate{
				Percentage: executorStatus.Percentage,
				Message:    executorStatus.Message,
				StatusType: statustype.Info,
			})
		},
	)

	err = parallelExecutor.Run(func(index *int, ctx *parallel.Context[fileutils.FileInfo]) {
		files := ctx.Items

		originCopy := reflect.New(reflect.TypeOf(origin).Elem()).Interface().(storage.Kind)
		err = originCopy.Initialize(origin.GetConfiguration())
		if err != nil {
			parallelExecutor.Error <- err
			return
		}

		destinationCopy := reflect.New(reflect.TypeOf(destination).Elem()).Interface().(storage.Kind)
		err := destinationCopy.Initialize(destination.GetConfiguration())
		if err != nil {
			parallelExecutor.Error <- err
			return
		}

		for _, file := range files {
			if file.Directory {
				continue
			}
			strategy.tryIncrementalCopy(originCopy, destinationCopy, file, backupId.String(), parallelExecutor.Error, blocks)

			if ctx.Data["proceededBytes"] == nil {
				ctx.Data["proceededBytes"] = uint64(0)
			}

			processedSize := ctx.Data["proceededBytes"].(uint64)
			ctx.Data["proceededBytes"] = processedSize + file.Size

			if sizes[*index] == 0 || ctx.Data["proceededBytes"].(uint64) == 0 {
				ctx.Status = 0
			} else {
				ctx.Status = float64(ctx.Data["proceededBytes"].(uint64)) / float64(sizes[*index])
			}
		}

		// Close the connections
		originCopy.Close()
		destinationCopy.Close()
		ctx.Finished = true
	})
	if err != nil {
		return nil, err
	}

	return strategy.handleIndexingBackup(backupId, jobId, origin, destination, callback)
}

func (strategy *IncrementalStrategy) handleIndexingBackup(id uuid.UUID, jobId uuid.UUID, origin storage.Kind,
	destination storage.Kind, callback func(StatusUpdate)) (*RawBackupInfo, error) {

	callback(StatusUpdate{
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
		callback(StatusUpdate{
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
		callback(StatusUpdate{
			Percentage: 0,
			Message:    err.Error(),
			StatusType: statustype.Error,
		})
		return nil, err
	}

	callback(StatusUpdate{
		Percentage: 100,
		Message:    "FINISHED",
		StatusType: statustype.Finished,
	})

	return &rawBackupInfo, nil
}

func (strategy *IncrementalStrategy) handleInitialize(destinationDirectory string, sourceDirectory string,
	from storage.Kind, to storage.Kind, callback func(StatusUpdate)) ([]fileutils.FileInfo, error) {
	callback(StatusUpdate{
		Percentage: 0,
		Message:    "Initializing strategy",
		StatusType: statustype.Info,
	})
	rawFiles, err := from.ListFiles(sourceDirectory)
	if err != nil {
		return nil, err
	}

	to.CreateDirectory(destinationDirectory)

	for _, file := range rawFiles {
		if file.Directory {
			to.CreateDirectory(destinationDirectory + file.RelativePath)

			if err != nil {
				callback(StatusUpdate{
					Percentage: 0,
					Message:    err.Error(),
					StatusType: statustype.Error,
				})
				return nil, err
			}
		}
	}

	callback(StatusUpdate{
		Percentage: 0,
		Message:    "Detected " + strconv.Itoa(len(rawFiles)) + " files",
		StatusType: statustype.Info,
	})

	if err != nil {
		callback(StatusUpdate{
			Percentage: 0,
			Message:    err.Error(),
			StatusType: statustype.Error,
		})
		return nil, err
	}

	return rawFiles, nil
}

func (strategy *IncrementalStrategy) tryIncrementalCopy(originCopy storage.Kind, destinationCopy storage.Kind, file fileutils.FileInfo,
	storePrefix string, parallelError chan error, blocks []string) error {
	retries := 0
	for {
		copyErr := strategy.copyIncremental(originCopy, destinationCopy, file, storePrefix, blocks)
		if copyErr == nil {
			break
		}
		retries++
		if retries > 15 {
			parallelError <- copyErr
			return copyErr
		}
		// Wait 1 second before retrying so that we don't overload the server
		time.Sleep(time.Second)
	}
	return nil
}

func (strategy *IncrementalStrategy) copyIncremental(origin storage.Kind, destination storage.Kind,
	file fileutils.FileInfo, storePrefix string, blocks []string) error {
	if file.Directory {
		return nil
	}

	// Read File Size
	size, err := origin.Size(file.RelativePath)
	if err != nil {
		return err
	}

	iterations := int64(size) / strategy.Configuration.BlockSize
	lastBlockSize := int64(size) % strategy.Configuration.BlockSize

	if lastBlockSize != 0 {
		iterations++
	}

	fileName := storePrefix + "/" + file.RelativePath
	if iterations == 0 {
		err = destination.Store(fileName, []byte{}, uint64(0))
		if err != nil {
			return err
		}
	}

	var blockLinksOfFile = make(map[int64]string)
	for index := int64(0); index < iterations; index++ {
		readSize := strategy.Configuration.BlockSize
		if index == iterations-1 && lastBlockSize != 0 {
			readSize = lastBlockSize
		}

		data, err := origin.Read(file.RelativePath, uint64(index*strategy.Configuration.BlockSize), uint64(readSize))
		if err != nil {
			return err
		}
		var hash = strategy.Hash(data)
		blockLinksOfFile[readSize] = hash

		if utils.Contains(blocks, hash) {
			continue
		}

		err = destination.Store(fileutils.FixedPathName(BlocksDirectory+hash), data, uint64(index*strategy.Configuration.BlockSize))
		if err != nil {
			return err
		}
	}
	data, err := json.Marshal(blockLinksOfFile)
	if err != nil {
		return err
	}
	err = destination.Store(fileName, data, 0)
	if err != nil {
		return err
	}
	return nil
}

func (strategy *IncrementalStrategy) Restore(origin storage.Kind, destination storage.Kind, id *uuid.UUID, callback func(StatusUpdate)) error {
	err := deleteAllFiles(origin, "", 1, callback)
	if err != nil {
		return err
	}

	rawFiles, err := strategy.handleInitialize("", id.String(), destination, origin, callback)
	if err != nil {
		return err
	}

	parallelFiles, sizes := fileutils.SplitArrayBySize(rawFiles, strategy.Configuration.Parallel)
	parallelExecutor := parallel.NewExecutor(
		parallelFiles,
		func(executorStatus parallel.ExecutorStatus) {
			callback(StatusUpdate{
				Percentage: executorStatus.Percentage,
				Message:    executorStatus.Message,
				StatusType: statustype.Info,
			})
		},
	)

	err = parallelExecutor.Run(func(index *int, ctx *parallel.Context[fileutils.FileInfo]) {
		files := ctx.Items

		originCopy := reflect.New(reflect.TypeOf(origin).Elem()).Interface().(storage.Kind)
		err = originCopy.Initialize(origin.GetConfiguration())
		if err != nil {
			parallelExecutor.Error <- err
			return
		}

		destinationCopy := reflect.New(reflect.TypeOf(destination).Elem()).Interface().(storage.Kind)
		err := destinationCopy.Initialize(destination.GetConfiguration())
		if err != nil {
			parallelExecutor.Error <- err
			return
		}

		for _, file := range files {
			if file.Directory {
				continue
			}
			strategy.tryRestoreCopy(destinationCopy, originCopy, file, "", id.String(), parallelExecutor.Error)

			if ctx.Data["proceededBytes"] == nil {
				ctx.Data["proceededBytes"] = uint64(0)
			}

			processedSize := ctx.Data["proceededBytes"].(uint64)
			ctx.Data["proceededBytes"] = processedSize + file.Size

			if sizes[*index] == 0 || ctx.Data["proceededBytes"].(uint64) == 0 {
				ctx.Status = 0
			} else {
				ctx.Status = float64(ctx.Data["proceededBytes"].(uint64)) / float64(sizes[*index])
			}
		}

		// Close the connections
		originCopy.Close()
		destinationCopy.Close()
		ctx.Finished = true
	})
	if err != nil {
		return err
	}

	callback(StatusUpdate{
		Percentage: 100,
		Message:    "FINISHED",
		StatusType: statustype.Finished,
	})
	destination.Close()
	origin.Close()

	return nil
}

func (strategy *IncrementalStrategy) tryRestoreCopy(originCopy storage.Kind, destinationCopy storage.Kind,
	file fileutils.FileInfo, storePrefix string, readPrefix string, parallelError chan error) error {
	retries := 0
	for {
		copyErr := strategy.copyRestore(originCopy, destinationCopy, file, storePrefix, readPrefix)
		if copyErr == nil {
			break
		}
		retries++
		if retries > 15 {
			parallelError <- copyErr
			return copyErr
		}
		// Wait 1 second before retrying so that we don't overload the server
		time.Sleep(time.Second)
	}
	return nil
}

func (strategy *IncrementalStrategy) copyRestore(origin storage.Kind, destination storage.Kind,
	file fileutils.FileInfo, storePrefix string, readPrefix string) error {
	if file.Directory {
		return nil
	}
	var blockLinksOfFile = make(map[int64]string)
	data, err := origin.Read(readPrefix+"/"+file.RelativePath, 0, 0)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &blockLinksOfFile)
	if err != nil {
		return err
	}
	cursorPosition := int64(0)

	fileName := storePrefix + "/" + file.RelativePath
	for size, hash := range blockLinksOfFile {
		data, err := origin.Read(fileutils.FixedPathName(BlocksDirectory+hash), uint64(0), uint64(size))
		if err != nil {
			return err
		}
		err = destination.Store(fileName, data, uint64(cursorPosition))
		if err != nil {
			return err
		}
		cursorPosition += size
	}
	return nil
}

func (strategy *IncrementalStrategy) Delete(destination storage.Kind, id *uuid.UUID, callback func(StatusUpdate)) error {
	err := destination.Delete(fileutils.FixedPathName(BackupInfoDirectory + id.String()))
	if err != nil {
		return err
	}
	deleteAllFiles(destination, id.String(), 1, func(status StatusUpdate) {
		callback(status)
	})
	destination.DeleteDirectory(id.String())

	callback(StatusUpdate{
		Percentage: 100,
		Message:    "FINISHED",
		StatusType: statustype.Finished,
	})

	return nil
}

func (strategy *IncrementalStrategy) Parse(configurationJson string) (interface{}, error) {
	var configuration IncrementalStrategyConfiguration
	err := json.Unmarshal([]byte(configurationJson), &configuration)
	if err != nil {
		return nil, err
	}
	return &configuration, nil
}

func (strategy *IncrementalStrategy) DefaultConfiguration() interface{} {
	return reflect.New(reflect.TypeOf(IncrementalStrategyConfiguration{})).Interface()
}

func (strategy *IncrementalStrategy) ListBackups(destination storage.Kind) ([]*RawBackupInfo, error) {
	handleCreateSlixxDirectories(destination)

	destination.CreateDirectory(fileutils.FixedPathName(BlocksDirectory))

	return nil, nil
}

func (strategy *IncrementalStrategy) Close() error {
	return nil
}
