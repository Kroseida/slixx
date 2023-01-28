package controller

import (
	"context"
	"github.com/google/uuid"
	"github.com/samsarahq/thunder/reactive"
	"kroseida.org/slixx/internal/master/application"
	"kroseida.org/slixx/internal/master/datasource"
	"kroseida.org/slixx/pkg/dto"
	"time"
)

type Storage struct {
	Id            uuid.UUID `json:"id",graphql:"id"`
	Name          string    `json:"name",graphql:"name"`
	Kind          string    `json:"kind",graphql:"kind"`
	Configuration string    `json:"configuration",graphql:"configuration"`
	CreatedAt     time.Time `json:"createdAt",graphql:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt",graphql:"updatedAt"`
	DeletedAt     time.Time `json:"deletedAt",graphql:"deletedAt"`
}

type StoragePrototype struct {
	Id        uuid.UUID `json:"id",graphql:"id"`
	Name      string    `json:"name",graphql:"name"`
	Kind      string    `json:"kind",graphql:"kind"`
	CreatedAt time.Time `json:"createdAt",graphql:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt",graphql:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt",graphql:"deletedAt"`
}

func GetStorage(ctx context.Context, args struct {
	Id uuid.UUID `json:"id",graphql:"id"`
}) (*Storage, error) {
	reactive.InvalidateAfter(ctx, 5*time.Second)
	storage, err := datasource.StorageProvider.GetStorage(args.Id)
	if err != nil {
		application.Logger.Debug(err)
		return nil, err
	}
	var storageDto Storage
	dto.Map(storage, &storageDto)

	return &storageDto, nil
}

func GetStorages(ctx context.Context) ([]*StoragePrototype, error) {
	reactive.InvalidateAfter(ctx, 5*time.Second)
	storages, err := datasource.StorageProvider.GetStorages()
	if err != nil {
		application.Logger.Debug(err)
		return nil, err
	}
	var storagesDto []*StoragePrototype
	dto.Map(storages, &storagesDto)

	return storagesDto, nil
}

func CreateStorage(ctx context.Context, args struct {
	Name          string `json:"name",graphql:"name"`
	Kind          string `json:"kind",graphql:"kind"`
	Configuration string `json:"configuration",graphql:"configuration"`
}) (*Storage, error) {
	reactive.InvalidateAfter(ctx, 5*time.Second)
	storage, err := datasource.StorageProvider.CreateStorage(args.Name, args.Kind, args.Configuration)
	if err != nil {
		application.Logger.Debug(err)
		return nil, err
	}
	var storageDto Storage
	dto.Map(storage, &storageDto)

	return &storageDto, err
}

func DeleteStorage(args struct {
	Id uuid.UUID `json:"id",graphql:"id"`
}) (bool, error) {
	result, err := datasource.StorageProvider.DeleteStorage(args.Id)
	if err != nil {
		application.Logger.Debug(err)
		return false, err
	}
	return result, nil
}
