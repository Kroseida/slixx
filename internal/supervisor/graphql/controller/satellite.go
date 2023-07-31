package controller

import (
	"context"
	"github.com/google/uuid"
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/reactive"
	"kroseida.org/slixx/internal/supervisor/datasource/provider"
	satelliteService "kroseida.org/slixx/internal/supervisor/service/satellite"
	"kroseida.org/slixx/pkg/dto"
	"kroseida.org/slixx/pkg/model"
	"time"
)

type Satellite struct {
	Id          uuid.UUID `sql:"default:uuid_generate_v4()"`
	Name        string
	Description string
	Address     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SatellitePrototype struct {
	Id          uuid.UUID `sql:"default:uuid_generate_v4()"`
	Name        string
	Description string
	Address     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type GetSatelliteDto struct {
	Id uuid.UUID `json:"id" graphql:"id"`
}

type SatellitesPage struct {
	Rows []SatellitePrototype `json:"rows" graphql:"rows"`
	Page
}

func GetSatellite(ctx context.Context, args GetSatelliteDto) (*Satellite, error) {
	if !IsPermitted(ctx, []string{"satellite.view"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	satellite, err := satelliteService.Get(args.Id)
	if err != nil {
		return nil, err
	}
	var satelliteDto *Satellite
	dto.Map(&satellite, &satelliteDto)

	return satelliteDto, nil
}

func GetSatellites(ctx context.Context, args PageArgs) (*SatellitesPage, error) {
	if !IsPermitted(ctx, []string{"satellite.view"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)

	var pagination provider.Pagination[model.Satellite]
	dto.Map(&args, &pagination)

	pages, err := satelliteService.GetPaged(&pagination)
	if err != nil {
		return nil, err
	}

	var pageDto SatellitesPage
	dto.Map(&pages, &pageDto)

	return &pageDto, nil
}

type CreateSatelliteDto struct {
	Name        string
	Description string
	Address     string
	Token       string
}

func CreateSatellite(ctx context.Context, args CreateSatelliteDto) (*Satellite, error) {
	if !IsPermitted(ctx, []string{"satellite.create"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)

	satellite, err := satelliteService.Create(args.Name, args.Description, args.Address, args.Token)

	if err != nil {
		return nil, err
	}
	var satellitesDto Satellite
	dto.Map(satellite, &satellitesDto)

	return &satellitesDto, err
}

type UpdateSatelliteDto struct {
	Id          uuid.UUID
	Name        *string
	Description *string
	Address     *string
	Token       *string
}

func UpdateSatellite(ctx context.Context, args UpdateSatelliteDto) (*Satellite, error) {
	if !IsPermitted(ctx, []string{"satellite.update"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	satellite, err := satelliteService.Update(
		args.Id,
		args.Name,
		args.Description,
		args.Address,
		args.Token,
	)
	if err != nil {
		return nil, err
	}
	var satelliteDto Satellite
	dto.Map(satellite, &satelliteDto)

	return &satelliteDto, err
}

type DeleteSatelliteDto struct {
	Id uuid.UUID `json:"id" graphql:"id"`
}

func DeleteSatellite(ctx context.Context, args DeleteSatelliteDto) (*Satellite, error) {
	if !IsPermitted(ctx, []string{"satellite.delete"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	satellite, err := satelliteService.Delete(args.Id)
	if err != nil {
		return nil, err
	}
	var satelliteDto Satellite
	dto.Map(satellite, &satelliteDto)

	return &satelliteDto, err
}
