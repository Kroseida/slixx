package strategy

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm/utils"
	"kroseida.org/slixx/pkg/statustype"
	"kroseida.org/slixx/pkg/storage"
	utils_ "kroseida.org/slixx/pkg/utils"
	"kroseida.org/slixx/pkg/utils/fileutils"
	"kroseida.org/slixx/pkg/utils/parallel"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

const BlocksDirectory = "blocks/"

type IncrementalStrategy struct {
	Configuration *IncrementalStrategyConfiguration
}

type IncrementalStrategyConfiguration struct {
	BlockSize int64 `json:"blockSize" slixx:"BYTE" default:"104857600"` // 100MB
	Parallel  int   `json:"parallel" slixx:"LONG" default:"4"`
}

type BlockReference struct {
	Hash   string `json:"hash"`
	Offset int64  `json:"offset"`
	Size   int64  `json:"size"`
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

func (strategy *IncrementalStrategy) Execute(job *parallel.RunningJob, origin storage.Kind, destination storage.Kind) (*RawBackupInfo, error) {
	if job.CheckCancel() {
		return nil, nil
	}
	handleCreateSlixxDirectories(destination)

	destination.CreateDirectory(fileutils.FixedPathName(BlocksDirectory))
	if job.CheckCancel() {
		return nil, nil
	}

	backupId := uuid.New()
	rawFiles, err := strategy.handleInitialize(job, backupId.String(), "", origin, destination)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	if job.CheckCancel() {
		return nil, nil
	}
	blockFiles, err := destination.ListFiles(BlocksDirectory)
	if err != nil {
		return nil, err
	}
	if job.CheckCancel() {
		return nil, nil
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
		job,
		parallelFiles,
		func(executorStatus parallel.ExecutorStatus) {
			job.Callback(utils_.StatusUpdate{
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
			strategy.tryIncrementalCopy(job, originCopy, destinationCopy, file, backupId.String(), parallelExecutor.Error, blocks)
			if job.CheckCancel() {
				return
			}

			if ctx.Data["proceededBytes"] == nil {
				ctx.Data["proceededBytes"] = uint64(0)
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
		return nil, err
	}
	if job.CheckCancel() {
		return nil, nil
	}

	return strategy.handleIndexingBackup(backupId, job, origin, destination)
}

func (strategy *IncrementalStrategy) handleIndexingBackup(id uuid.UUID, job *parallel.RunningJob, origin storage.Kind,
	destination storage.Kind) (*RawBackupInfo, error) {
	if job.CheckCancel() {
		return nil, nil
	}

	job.Callback(utils_.StatusUpdate{
		Percentage: 100,
		Message:    "Creating backup info file on destination",
		StatusType: statustype.Info,
	})

	handleCreateSlixxDirectories(destination)
	if job.CheckCancel() {
		return nil, nil
	}

	// This is the raw backup info that we will store in the destination storage, so that we can restore it later e.g
	// after supervisor get corrupted data or something like that we can restore this backup info and continue
	rawBackupInfo := RawBackupInfo{
		Id:              &id,
		CreatedAt:       time.Now(),
		JobId:           &job.JobId,
		OriginKind:      origin.GetName(),      // Store origin kind so that we can restore it later
		DestinationKind: destination.GetName(), // Store destination kind so that we can restore it later
		Strategy:        strategy.GetName(),    // Store strategy so that we can restore it later
	}
	rawBackupInfoBytes, err := json.Marshal(rawBackupInfo)
	if err != nil {
		job.Callback(utils_.StatusUpdate{
			Percentage: 0,
			Message:    err.Error(),
			StatusType: statustype.Error,
		})
		return nil, err
	}
	if job.CheckCancel() {
		return nil, nil
	}

	// Store backup info file in destination so that we have some information about the backup
	err = destination.Store(
		fileutils.FixedPathName(BackupInfoDirectory+"/"+id.String()),
		rawBackupInfoBytes,
		0,
	)
	if err != nil {
		job.Callback(utils_.StatusUpdate{
			Percentage: 0,
			Message:    err.Error(),
			StatusType: statustype.Error,
		})
		return nil, err
	}

	job.Callback(utils_.StatusUpdate{
		Percentage: 100,
		Message:    "FINISHED",
		StatusType: statustype.Finished,
	})

	return &rawBackupInfo, nil
}

func (strategy *IncrementalStrategy) handleInitialize(job *parallel.RunningJob, destinationDirectory string, sourceDirectory string,
	from storage.Kind, to storage.Kind) ([]fileutils.FileInfo, error) {
	if job.CheckCancel() {
		return nil, nil
	}

	job.Callback(utils_.StatusUpdate{
		Percentage: 0,
		Message:    "Initializing strategy",
		StatusType: statustype.Info,
	})
	rawFiles, err := from.ListFiles(sourceDirectory)
	if err != nil {
		return nil, err
	}

	to.CreateDirectory(destinationDirectory)
	if job.CheckCancel() {
		return nil, nil
	}

	for _, file := range rawFiles {
		if file.Directory {
			to.CreateDirectory(destinationDirectory + file.RelativePath)

			if err != nil {
				job.Callback(utils_.StatusUpdate{
					Percentage: 0,
					Message:    err.Error(),
					StatusType: statustype.Error,
				})
				return nil, err
			}
		}
	}
	if job.CheckCancel() {
		return nil, nil
	}

	job.Callback(utils_.StatusUpdate{
		Percentage: 0,
		Message:    "Detected " + strconv.Itoa(len(rawFiles)) + " files",
		StatusType: statustype.Info,
	})

	if err != nil {
		job.Callback(utils_.StatusUpdate{
			Percentage: 0,
			Message:    err.Error(),
			StatusType: statustype.Error,
		})
		return nil, err
	}
	if job.CheckCancel() {
		return nil, nil
	}

	return rawFiles, nil
}

func (strategy *IncrementalStrategy) tryIncrementalCopy(job *parallel.RunningJob, originCopy storage.Kind, destinationCopy storage.Kind, file fileutils.FileInfo,
	storePrefix string, parallelError chan error, blocks []string) error {
	retries := 0
	for {
		if job.CheckCancel() {
			return nil
		}

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

	var blockLinksOfFile = make([]*BlockReference, iterations)
	for index := int64(0); index < iterations; index++ {
		readSize := strategy.Configuration.BlockSize
		if index == iterations-1 && lastBlockSize != 0 {
			readSize = lastBlockSize
		}

		data, err := origin.Read(file.RelativePath, uint64(index*strategy.Configuration.BlockSize), uint64(readSize))

		if err != nil {
			return err
		}
		// Without last 2 characters
		var hash = strategy.Hash(data)
		blockLinksOfFile[index] = &BlockReference{
			Hash:   hash,
			Offset: index * strategy.Configuration.BlockSize,
			Size:   readSize,
		}

		if utils.Contains(blocks, hash) {
			continue
		}
		err = destination.Store(fileutils.FixedPathName(BlocksDirectory+"/"+hash), data, 0)
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

func (strategy *IncrementalStrategy) Restore(job *parallel.RunningJob, origin storage.Kind, destination storage.Kind, id *uuid.UUID) error {
	err := deleteAllFiles(job, origin, "", 1)
	if err != nil {
		return err
	}

	rawFiles, err := strategy.handleInitialize(job, "", id.String(), destination, origin)
	if err != nil {
		return err
	}
	if job.CheckCancel() {
		return nil
	}
	parallelFiles, sizes := fileutils.SplitArrayBySize(rawFiles, strategy.Configuration.Parallel)
	parallelExecutor := parallel.NewExecutor(
		job,
		parallelFiles,
		func(executorStatus parallel.ExecutorStatus) {
			job.Callback(utils_.StatusUpdate{
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

		destinationCopy := reflect.New(reflect.TypeOf(destination).Elem()).Interface().(storage.Kind)
		err := destinationCopy.Initialize(destination.GetConfiguration())
		if err != nil {
			parallelExecutor.Error <- err
			return
		}

		for _, file := range files {
			if job.CheckCancel() {
				return
			}

			if file.Directory {
				continue
			}
			strategy.tryRestoreCopy(job, destinationCopy, originCopy, file, "", id.String(), parallelExecutor.Error)
			if job.CheckCancel() {
				return
			}

			if ctx.Data["proceededBytes"] == nil {
				ctx.Data["proceededBytes"] = uint64(0)
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

	job.Callback(utils_.StatusUpdate{
		Percentage: 100,
		Message:    "FINISHED",
		StatusType: statustype.Finished,
	})
	destination.Close()
	origin.Close()

	return nil
}

func (strategy *IncrementalStrategy) tryRestoreCopy(job *parallel.RunningJob, originCopy storage.Kind, destinationCopy storage.Kind,
	file fileutils.FileInfo, storePrefix string, readPrefix string, parallelError chan error) error {
	retries := 0
	for {
		if job.CheckCancel() {
			return nil
		}
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
	var blockLinksOfFile = make([]*BlockReference, 0)
	sort.Slice(blockLinksOfFile, func(i, j int) bool {
		return blockLinksOfFile[i].Offset > blockLinksOfFile[j].Offset
	})

	data, err := origin.Read(readPrefix+"/"+file.RelativePath, 0, 0)
	if err != nil {
		return errors.New("error while loading block links for file: " + file.RelativePath + " error: " + err.Error())
	}
	err = json.Unmarshal(data, &blockLinksOfFile)
	if err != nil {
		return errors.New("error while loading block links for file: " + file.RelativePath + " error: " + err.Error())
	}

	fileName := fileutils.FixedPathName(storePrefix + "/" + file.RelativePath)
	for _, reference := range blockLinksOfFile {
		data, err := origin.Read(fileutils.FixedPathName(BlocksDirectory+reference.Hash), uint64(0), uint64(reference.Size))
		if err != nil {
			return errors.New("error while reading block: " + reference.Hash + " for file: " + fileName + " error: " + err.Error())
		}
		err = destination.Store(fileName, data, uint64(reference.Offset))
		if err != nil {
			return errors.New("error while storing block: " + reference.Hash + " for file: " + fileName + " error: " + err.Error())
		}
	}
	return nil
}

func (strategy *IncrementalStrategy) Delete(job *parallel.RunningJob, destination storage.Kind, id *uuid.UUID) error {
	err := destination.Delete(fileutils.FixedPathName(BackupInfoDirectory + id.String()))
	if err != nil {
		return err
	}
	deleteAllFiles(job, destination, id.String(), 1)
	destination.DeleteDirectory(id.String())

	job.Callback(utils_.StatusUpdate{
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

	files, err := destination.ListFiles(BackupInfoDirectory)
	if err != nil {
		return nil, err
	}

	backupList := make([]*RawBackupInfo, 0, len(files))
	for _, file := range files {
		idAsString := strings.TrimPrefix(file.RelativePath, BackupInfoDirectory)
		if idAsString[0] == '/' || idAsString[0] == '\\' {
			idAsString = idAsString[1:]
		}

		parsedId, err := uuid.Parse(idAsString)
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

func (strategy *IncrementalStrategy) Close() error {
	return nil
}
