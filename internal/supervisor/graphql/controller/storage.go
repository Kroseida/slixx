package controller

import (
	"context"
	"github.com/google/uuid"
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/reactive"
	"kroseida.org/slixx/internal/supervisor/application"
	"kroseida.org/slixx/internal/supervisor/datasource/provider"
	storageService "kroseida.org/slixx/internal/supervisor/service/storage"
	"kroseida.org/slixx/pkg/dto"
	"kroseida.org/slixx/pkg/model"
	"kroseida.org/slixx/pkg/storage"
	"reflect"
	"time"
)

type Storage struct {
	Id            uuid.UUID `json:"id" graphql:"id"`
	Name          string    `json:"name" graphql:"name"`
	Description   string    `json:"description" graphql:"description"`
	Kind          string    `json:"kind" graphql:"kind"`
	Configuration string    `json:"configuration" graphql:"configuration"`
	CreatedAt     time.Time `json:"createdAt" graphql:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt" graphql:"updatedAt"`
	DeletedAt     time.Time `json:"deletedAt" graphql:"deletedAt"`
}

type StoragePrototype struct {
	Id          uuid.UUID `json:"id" graphql:"id"`
	Name        string    `json:"name" graphql:"name"`
	Description string    `json:"description" graphql:"description"`
	Kind        string    `json:"kind" graphql:"kind"`
	CreatedAt   time.Time `json:"createdAt" graphql:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" graphql:"updatedAt"`
	DeletedAt   time.Time `json:"deletedAt" graphql:"deletedAt"`
}

type StoragesPage struct {
	Rows []StoragePrototype `json:"rows" graphql:"rows"`
	Page
}

type GetStorageDto struct {
	Id uuid.UUID `json:"id" graphql:"id"`
}

func GetStorage(ctx context.Context, args GetStorageDto) (*Storage, error) {
	if !IsPermitted(ctx, []string{"storage.view"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	storage, err := storageService.Get(args.Id)
	if err != nil {
		application.Logger.Debug(err)
		return nil, err
	}
	var storageDto *Storage
	dto.Map(&storage, &storageDto)

	return storageDto, nil
}

func GetStorages(ctx context.Context, args PageArgs) (*StoragesPage, error) {
	if !IsPermitted(ctx, []string{"storage.view"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)

	var pagination provider.Pagination[model.Storage]
	dto.Map(&args, &pagination)

	pages, err := storageService.GetPaged(&pagination)
	if err != nil {
		application.Logger.Debug(err)
		return nil, err
	}

	var pageDtos StoragesPage
	dto.Map(&pages, &pageDtos)

	return &pageDtos, nil
}

type CreateStorageDto struct {
	Name          string `json:"name" graphql:"name"`
	Description   string `json:"description" graphql:"description"`
	Kind          string `json:"kind" graphql:"kind"`
	Configuration string `json:"configuration" graphql:"configuration"`
}

func CreateStorage(ctx context.Context, args CreateStorageDto) (*Storage, error) {
	if !IsPermitted(ctx, []string{"storage.create"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	storage, err := storageService.Create(args.Name, args.Description, args.Kind, args.Configuration)
	if err != nil {
		application.Logger.Debug(err)
		return nil, err
	}
	var storageDto Storage
	dto.Map(storage, &storageDto)

	return &storageDto, err
}

type UpdateStorageDto struct {
	Id            uuid.UUID `json:"id" graphql:"id"`
	Name          *string   `json:"name" graphql:"name"`
	Description   *string   `json:"description" graphql:"description"`
	Kind          *string   `json:"kind" graphql:"kind"`
	Configuration *string   `json:"configuration" graphql:"configuration"`
}

func UpdateStorage(ctx context.Context, args UpdateStorageDto) (*Storage, error) {
	if !IsPermitted(ctx, []string{"storage.update"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	storage, err := storageService.Update(
		args.Id,
		args.Name,
		args.Description,
		args.Kind,
		args.Configuration,
	)

	if err != nil {
		application.Logger.Debug(err)
		return nil, err
	}
	var storageDto Storage
	dto.Map(&storage, &storageDto)

	return &storageDto, nil
}

type DeleteStorageDto struct {
	Id uuid.UUID `json:"id" graphql:"id"`
}

func DeleteStorage(ctx context.Context, args DeleteStorageDto) (*Storage, error) {
	if !IsPermitted(ctx, []string{"storage.delete"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	storage, err := storageService.Delete(args.Id)
	if err != nil {
		application.Logger.Debug(err)
		return nil, err
	}
	var storageDto Storage
	dto.Map(&storage, &storageDto)

	return &storageDto, nil
}

type StorageKindDescriptionDto struct {
	Name          string                               `json:"name" graphql:"name"`
	Configuration []StorageConfigurationDescriptionDto `json:"configuration" graphql:"configuration"`
}

type StorageConfigurationDescriptionDto struct {
	Name    string `json:"name" graphql:"name"`
	Kind    string `json:"kind" graphql:"kind"`
	Default string `json:"default" graphql:"default"`
}

func GetStorageKinds() ([]StorageKindDescriptionDto, error) {
	var descriptions []StorageKindDescriptionDto
	for _, kind := range storage.Values() {
		var configurations []StorageConfigurationDescriptionDto

		val := reflect.ValueOf(kind.DefaultConfiguration()).Elem()
		for i := 0; i < val.NumField(); i++ {
			configurations = append(configurations, StorageConfigurationDescriptionDto{
				Name:    val.Type().Field(i).Tag.Get("json"),
				Kind:    val.Type().Field(i).Tag.Get("slixx"),
				Default: val.Type().Field(i).Tag.Get("default"),
			})
		}

		descriptions = append(descriptions, StorageKindDescriptionDto{
			Name:          kind.GetName(),
			Configuration: configurations,
		})
	}
	return descriptions, nil
}
