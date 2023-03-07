package controller

import (
	"context"
	"github.com/google/uuid"
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/reactive"
	"kroseida.org/slixx/internal/master/application"
	"kroseida.org/slixx/internal/master/datasource"
	"kroseida.org/slixx/internal/master/datasource/provider"
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

func GetStorage(ctx context.Context, args struct {
	Id uuid.UUID `json:"id" graphql:"id"`
}) (*Storage, error) {
	if !IsPermitted(ctx, []string{"storage.view"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
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

func GetStorages(ctx context.Context, args PageArgs) (*StoragesPage, error) {
	if !IsPermitted(ctx, []string{"storage.view"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)

	var pagination provider.Pagination[model.Storage]
	dto.Map(&args, &pagination)

	pages, err := datasource.StorageProvider.GetStoragesPaged(&pagination)
	if err != nil {
		application.Logger.Debug(err)
		return nil, err
	}

	var pageDtos StoragesPage
	dto.Map(&pages, &pageDtos)

	return &pageDtos, nil
}

func CreateStorage(ctx context.Context, args struct {
	Name          string `json:"name" graphql:"name"`
	Description   string `json:"description" graphql:"description"`
	Kind          string `json:"kind" graphql:"kind"`
	Configuration string `json:"configuration" graphql:"configuration"`
}) (*Storage, error) {
	if !IsPermitted(ctx, []string{"storage.create"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	storage, err := datasource.StorageProvider.CreateStorage(args.Name, args.Description, args.Kind, args.Configuration)
	if err != nil {
		application.Logger.Debug(err)
		return nil, err
	}
	var storageDto Storage
	dto.Map(storage, &storageDto)

	return &storageDto, err
}

func UpdateStorage(ctx context.Context, args struct {
	Id            uuid.UUID `json:"id" graphql:"id"`
	Name          *string   `json:"name" graphql:"name"`
	Description   *string   `json:"description" graphql:"description"`
	Kind          *string   `json:"kind" graphql:"kind"`
	Configuration *string   `json:"configuration" graphql:"configuration"`
}) (*Storage, error) {
	if !IsPermitted(ctx, []string{"storage.update"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	storage, err := datasource.StorageProvider.UpdateStorage(args.Id, args.Name, args.Description, args.Kind, args.Configuration)

	if err != nil {
		application.Logger.Debug(err)
		return nil, err
	}
	var storageDto Storage
	dto.Map(&storage, &storageDto)

	return &storageDto, nil
}

func DeleteStorage(ctx context.Context, args struct {
	Id uuid.UUID `json:"id" graphql:"id"`
}) (*Storage, error) {
	if !IsPermitted(ctx, []string{"storage.delete"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	storage, err := datasource.StorageProvider.DeleteStorage(args.Id)
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
	Name string `json:"name" graphql:"name"`
	Kind string `json:"kind" graphql:"kind"`
}

func GetStorageKinds() ([]StorageKindDescriptionDto, error) {
	var descriptions []StorageKindDescriptionDto
	for _, kind := range storage.Values() {
		var configurations []StorageConfigurationDescriptionDto

		val := reflect.ValueOf(kind.DefaultConfiguration()).Elem()
		for i := 0; i < val.NumField(); i++ {
			configurations = append(configurations, StorageConfigurationDescriptionDto{
				Name: val.Type().Field(i).Tag.Get("json"),
				Kind: val.Type().Field(i).Tag.Get("slixx"),
			})
		}

		descriptions = append(descriptions, StorageKindDescriptionDto{
			Name:          kind.GetName(),
			Configuration: configurations,
		})
	}

	return descriptions, nil
}
