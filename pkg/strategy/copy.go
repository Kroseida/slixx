package strategy

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"kroseida.org/slixx/pkg/statustype"
	"kroseida.org/slixx/pkg/storage"
	"kroseida.org/slixx/pkg/utils"
	"kroseida.org/slixx/pkg/utils/fileutils"
	"kroseida.org/slixx/pkg/utils/parallel"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const SlixxDirectory = ".slixx"
const BackupInfoDirectory = SlixxDirectory + "/backups/"

type CopyStrategy struct {
	Configuration *CopyStrategyConfiguration
}

type CopyStrategyConfiguration struct {
	BlockSize int `json:"blockSize" slixx:"BYTE" default:"1073741824"` // 1GB in bytes
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

func (strategy *CopyStrategy) Execute(job *parallel.RunningJob, origin storage.Kind, destination storage.Kind) (*RawBackupInfo, error) {
	backupId := uuid.New()
	rawFiles, err := strategy.handleInitialize(job, backupId.String(), "", origin, destination)
	if err != nil {
		return nil, err
	}
	parallelFiles, sizes := fileutils.SplitArrayBySize(rawFiles, strategy.Configuration.Parallel)
	if job.CheckCancel() {
		return nil, nil
	}

	parallelExecutor := parallel.NewExecutor(
		job,
		parallelFiles,
		func(executorStatus parallel.ExecutorStatus) {
			job.Callback(utils.StatusUpdate{
				Percentage: executorStatus.Percentage,
				Message:    executorStatus.Message,
				StatusType: statustype.Info,
			})
		},
	)

	if job.CheckCancel() {
		return nil, nil
	}

	err = parallelExecutor.Run(func(index *int, ctx *parallel.Context[fileutils.FileInfo]) {
		files := ctx.Items

		if job.CheckCancel() {
			return
		}

		originCopy := reflect.New(reflect.TypeOf(origin).Elem()).Interface().(storage.Kind)
		err = originCopy.Initialize(origin.GetConfiguration())
		if err != nil {
			parallelExecutor.Error <- err
			return
		}
		if job.CheckCancel() {
			return
		}

		destinationCopy := reflect.New(reflect.TypeOf(destination).Elem()).Interface().(storage.Kind)
		err := destinationCopy.Initialize(destination.GetConfiguration())
		if err != nil {
			parallelExecutor.Error <- err
			return
		}
		if job.CheckCancel() {
			return
		}

		for _, file := range files {
			if job.CheckCancel() {
				return
			}

			if file.Directory {
				continue
			}
			strategy.tryCopy(job, originCopy, destinationCopy, file, backupId.String(), "", parallelExecutor.Error)
			if job.CheckCancel() {
				return
			}

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

	return strategy.handleIndexingBackup(job, backupId, origin, destination)
}

func (strategy *CopyStrategy) tryCopy(job *parallel.RunningJob, originCopy storage.Kind, destinationCopy storage.Kind, file fileutils.FileInfo,
	storePrefix string, readPrefix string, parallelError chan error) error {
	retries := 0
	for {
		if job.CheckCancel() {
			return nil
		}
		copyErr := strategy.copy(originCopy, destinationCopy, file, storePrefix, readPrefix)
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

func (strategy *CopyStrategy) handleInitialize(job *parallel.RunningJob, destinationDirectory string, sourceDirectory string, from storage.Kind, to storage.Kind) ([]fileutils.FileInfo, error) {
	job.Callback(utils.StatusUpdate{
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
			err := to.CreateDirectory(destinationDirectory + file.RelativePath)
			if err != nil {
				job.Callback(utils.StatusUpdate{
					Percentage: 0,
					Message:    err.Error(),
					StatusType: statustype.Error,
				})
				return nil, err
			}
		}
	}

	job.Callback(utils.StatusUpdate{
		Percentage: 0,
		Message:    "Detected " + strconv.Itoa(len(rawFiles)) + " files",
		StatusType: statustype.Info,
	})

	if err != nil {
		job.Callback(utils.StatusUpdate{
			Percentage: 0,
			Message:    err.Error(),
			StatusType: statustype.Error,
		})
		return nil, err
	}

	return rawFiles, nil
}

func (strategy *CopyStrategy) handleIndexingBackup(job *parallel.RunningJob, id uuid.UUID, origin storage.Kind, destination storage.Kind) (*RawBackupInfo, error) {

	job.Callback(utils.StatusUpdate{
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
		JobId:           &job.JobId,            // Store job id so that we can restore it later
		OriginKind:      origin.GetName(),      // Store origin kind so that we can restore it later
		DestinationKind: destination.GetName(), // Store destination kind so that we can restore it later
		Strategy:        strategy.GetName(),    // Store strategy so that we can restore it later
	}
	rawBackupInfoBytes, err := json.Marshal(rawBackupInfo)
	if err != nil {
		job.Callback(utils.StatusUpdate{
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
		job.Callback(utils.StatusUpdate{
			Percentage: 0,
			Message:    err.Error(),
			StatusType: statustype.Error,
		})
		return nil, err
	}

	job.Callback(utils.StatusUpdate{
		Percentage: 100,
		Message:    "FINISHED",
		StatusType: statustype.Finished,
	})

	return &rawBackupInfo, nil
}

func handleCreateSlixxDirectories(destination storage.Kind) {
	destination.CreateDirectory(fileutils.FixedPathName(SlixxDirectory))
	destination.CreateDirectory(fileutils.FixedPathName(BackupInfoDirectory))
}

func deleteAllFiles(job *parallel.RunningJob, targetStorage storage.Kind, prefix string, parallels int) error {
	rawFiles, err := targetStorage.ListFiles(prefix)
	if err != nil {
		return err
	}

	parallelFiles, sizes := fileutils.SplitArrayBySize(rawFiles, parallels)

	parallelExecutor := parallel.NewExecutor(
		job,
		parallelFiles,
		func(executorStatus parallel.ExecutorStatus) {
			job.Callback(utils.StatusUpdate{
				Percentage: executorStatus.Percentage,
				Message:    "Deleting files...",
				StatusType: statustype.Info,
			})
		},
	)

	parallelExecutor.Run(func(index *int, ctx *parallel.Context[fileutils.FileInfo]) {
		files := fileutils.SortByLength(ctx.Items)

		storageCopy := reflect.New(reflect.TypeOf(targetStorage).Elem()).Interface().(storage.Kind)
		err = storageCopy.Initialize(targetStorage.GetConfiguration())
		if err != nil {
			parallelExecutor.Error <- err
			return
		}

		for _, file := range files {
			if file.Directory {
				continue
			}
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

			err := storageCopy.Delete(prefix + "/" + file.RelativePath)
			if err != nil {
				parallelExecutor.Error <- err
				return
			}
		}

		for _, file := range files {
			if !file.Directory {
				continue
			}
			err := storageCopy.Delete(file.FullDirectory)
			if err != nil {
				parallelExecutor.Error <- err
				return
			}
		}

		storageCopy.Close()
		ctx.Finished = true
	})

	// Try to delete the files 10 times
	for retryDelete := 0; retryDelete < 10; retryDelete++ {
		// Delete again
		rawFiles, err = targetStorage.ListFiles(prefix)
		if err != nil {
			job.Callback(utils.StatusUpdate{
				Percentage: 0,
				Message:    err.Error(),
				StatusType: statustype.Error,
			})
		}

		// If there are no files left, then we are done
		if len(rawFiles) == 0 {
			break
		}

		for _, file := range fileutils.SortByLength(rawFiles) {
			targetStorage.DeleteDirectory(file.RelativePath)
		}
	}
	return nil
}

func (strategy *CopyStrategy) Delete(job *parallel.RunningJob, destination storage.Kind, id *uuid.UUID) error {
	err := destination.Delete(fileutils.FixedPathName(BackupInfoDirectory + id.String()))
	if err != nil {
		return err
	}
	deleteAllFiles(job, destination, id.String(), strategy.Configuration.Parallel)
	destination.DeleteDirectory(id.String())

	job.Callback(utils.StatusUpdate{
		Percentage: 100,
		Message:    "FINISHED",
		StatusType: statustype.Finished,
	})

	return nil
}

func (strategy *CopyStrategy) Restore(job *parallel.RunningJob, origin storage.Kind, destination storage.Kind, id *uuid.UUID) error {
	if job.CheckCancel() {
		return nil
	}

	err := deleteAllFiles(job, origin, "", strategy.Configuration.Parallel)
	if err != nil {
		return err
	}
	if job.CheckCancel() {
		return nil
	}

	rawFiles, err := strategy.handleInitialize(job, "", id.String(), destination, origin)
	if err != nil {
		return err
	}
	parallelFiles, sizes := fileutils.SplitArrayBySize(rawFiles, strategy.Configuration.Parallel)
	if job.CheckCancel() {
		return nil
	}

	parallelExecutor := parallel.NewExecutor(
		job,
		parallelFiles,
		func(executorStatus parallel.ExecutorStatus) {
			job.Callback(utils.StatusUpdate{
				Percentage: executorStatus.Percentage,
				Message:    executorStatus.Message,
				StatusType: statustype.Info,
			})
		},
	)

	err = parallelExecutor.Run(func(index *int, ctx *parallel.Context[fileutils.FileInfo]) {
		if job.CheckCancel() {
			return
		}

		files := ctx.Items

		originCopy := reflect.New(reflect.TypeOf(origin).Elem()).Interface().(storage.Kind)
		err = originCopy.Initialize(origin.GetConfiguration())
		if err != nil {
			parallelExecutor.Error <- err
			return
		}
		if job.CheckCancel() {
			return
		}

		destinationCopy := reflect.New(reflect.TypeOf(destination).Elem()).Interface().(storage.Kind)
		err := destinationCopy.Initialize(destination.GetConfiguration())
		if err != nil {
			parallelExecutor.Error <- err
			return
		}
		if job.CheckCancel() {
			return
		}

		for _, file := range files {
			if job.CheckCancel() {
				return
			}

			if file.Directory {
				continue
			}
			strategy.tryCopy(job, destinationCopy, originCopy, file, "", id.String(), parallelExecutor.Error)
			if job.CheckCancel() {
				return
			}

			if ctx.Data["proceededBytes"] == nil {
				ctx.Data["proceededBytes"] = uint64(0)
			}
			if job.CheckCancel() {
				return
			}

			processedSize := ctx.Data["proceededBytes"].(uint64)
			ctx.Data["proceededBytes"] = processedSize + file.Size
			if job.CheckCancel() {
				return
			}

			if sizes[*index] == 0 || ctx.Data["proceededBytes"].(uint64) == 0 {
				ctx.Status = 0
			} else {
				ctx.Status = float64(ctx.Data["proceededBytes"].(uint64)) / float64(sizes[*index])
			}
			if job.CheckCancel() {
				return
			}
		}
		if job.CheckCancel() {
			return
		}

		// Close the connections
		originCopy.Close()
		destinationCopy.Close()
		ctx.Finished = true
	})
	if err != nil {
		return err
	}
	if job.CheckCancel() {
		return nil
	}

	job.Callback(utils.StatusUpdate{
		Percentage: 100,
		Message:    "FINISHED",
		StatusType: statustype.Finished,
	})
	destination.Close()
	origin.Close()

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

	destination.Close()
	return backupList, nil
}

func (strategy *CopyStrategy) copy(origin storage.Kind, destination storage.Kind, file fileutils.FileInfo, storePrefix string, readPrefix string) error {
	if file.Directory {
		return nil
	}

	// Read File Size
	size, err := origin.Size(readPrefix + "/" + file.RelativePath)
	if err != nil {
		return errors.New("error while reading file size for " + file.RelativePath + ": " + err.Error())
	}

	iterations := int(size) / strategy.Configuration.BlockSize
	lastBlockSize := int(size) % strategy.Configuration.BlockSize

	if lastBlockSize != 0 {
		iterations++
	}

	fileName := storePrefix + "/" + file.RelativePath
	if iterations == 0 {
		err = destination.Store(fileName, []byte{}, uint64(0))
		if err != nil {
			return errors.New("error while storing empty file " + fileName + ": " + err.Error())
		}
	}

	for index := 0; index < iterations; index++ {
		readSize := strategy.Configuration.BlockSize
		if index == iterations-1 && lastBlockSize != 0 {
			readSize = lastBlockSize
		}

		data, err := origin.Read(readPrefix+"/"+file.RelativePath, uint64(index*strategy.Configuration.BlockSize), uint64(readSize))
		if err != nil {
			return errors.New("error while reading file " + readPrefix + "/" + file.RelativePath + ": " + err.Error())
		}

		err = destination.Store(fileName, data, uint64(index*strategy.Configuration.BlockSize))
		if err != nil {
			return errors.New("error while storing file " + readPrefix + "/" + fileName + ": " + err.Error())
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
