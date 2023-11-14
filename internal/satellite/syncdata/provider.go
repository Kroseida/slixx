package syncdata

import (
	"encoding/json"
	"github.com/google/uuid"
	"kroseida.org/slixx/internal/common"
	"kroseida.org/slixx/internal/satellite/application"
	"kroseida.org/slixx/pkg/model"
	"os"
)

type ContainerModel struct {
	Version  string
	Storages map[uuid.UUID]*model.Storage
	Jobs     map[uuid.UUID]*model.Job
}

var CacheFile = ".cache"
var DefaultContainer = ContainerModel{
	Version:  common.CurrentVersion,
	Storages: map[uuid.UUID]*model.Storage{},
	Jobs:     map[uuid.UUID]*model.Job{},
}
var Container = DefaultContainer

func GenerateCache() error {
	cacheJson, err := json.Marshal(&Container)
	if err != nil {
		return err
	}

	err = os.WriteFile(CacheFile, cacheJson, 0644)
	if err != nil {
		return err
	}

	return nil
}

func LoadCache() error {
	cacheJson, err := os.ReadFile(CacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			GenerateCache()
			return nil
		}
		return err
	}

	err = json.Unmarshal(cacheJson, &Container)
	if err != nil {
		return err
	}

	if Container.Version != common.CurrentVersion {
		Container = DefaultContainer
		application.Logger.Info("Cache version mismatch, renaming cache file to .cache.old")
		err := os.Rename(".cache", ".cache.old")
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}
