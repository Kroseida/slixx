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

type SatelliteDto struct {
	Id          uuid.UUID `sql:"default:uuid_generate_v4()"`
	Name        string
	Description string
	Address     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Connected   bool
}

type SatellitePrototypeDto struct {
	Id          uuid.UUID `sql:"default:uuid_generate_v4()"`
	Name        string
	Description string
	Address     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Connected   bool
}

type GetSatelliteDto struct {
	Id uuid.UUID `json:"id" graphql:"id"`
}

type SatellitesPageDto struct {
	Rows []SatellitePrototypeDto `json:"rows" graphql:"rows"`
	Page
}

func GetSatellite(ctx context.Context, args GetSatelliteDto) (*SatelliteDto, error) {
	if !IsPermitted(ctx, []string{"satellite.view"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	satellite, err := satelliteService.Get(args.Id)
	if err != nil {
		return nil, err
	}
	var satelliteDto *SatelliteDto
	dto.Map(&satellite, &satelliteDto)

	return satelliteDto, nil
}

func GetSatellites(ctx context.Context, args GetPageDto) (*SatellitesPageDto, error) {
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

	var pageDto SatellitesPageDto
	dto.Map(&pages, &pageDto)

	return &pageDto, nil
}

type CreateSatelliteDto struct {
	Name        string
	Description string
	Address     string
	Token       string
}

func CreateSatellite(ctx context.Context, args CreateSatelliteDto) (*SatelliteDto, error) {
	if !IsPermitted(ctx, []string{"satellite.create"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)

	satellite, err := satelliteService.Create(args.Name, args.Description, args.Address, args.Token)

	if err != nil {
		return nil, err
	}
	var satellitesDto SatelliteDto
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

func UpdateSatellite(ctx context.Context, args UpdateSatelliteDto) (*SatelliteDto, error) {
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
	var satelliteDto SatelliteDto
	dto.Map(satellite, &satelliteDto)

	return &satelliteDto, err
}

type DeleteSatelliteDto struct {
	Id uuid.UUID `json:"id" graphql:"id"`
}

func DeleteSatellite(ctx context.Context, args DeleteSatelliteDto) (*SatelliteDto, error) {
	if !IsPermitted(ctx, []string{"satellite.delete"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	satellite, err := satelliteService.Delete(args.Id)
	if err != nil {
		return nil, err
	}
	var satelliteDto SatelliteDto
	dto.Map(satellite, &satelliteDto)

	return &satelliteDto, err
}

type LogEntryDto struct {
	Id          uuid.UUID
	Sender      string
	SatelliteId uuid.UUID
	Level       string
	Message     string
	LoggedAt    time.Time
}

type GetLogsRequestDto struct {
	SatelliteId uuid.UUID `json:"satelliteId" graphql:"satelliteId"`
	Limit       *int64    `json:"limit,omitempty;query:limit"`
	Page        *int64    `json:"page,omitempty;query:page"`
	Sort        *string   `json:"sort,omitempty;query:sort"`
	Search      *string   `json:"search"`
}

type LogsPageDto struct {
	Rows []LogEntryDto `json:"rows" graphql:"rows"`
	Page
}

func GetSatelliteLogs(ctx context.Context, args GetLogsRequestDto) (*LogsPageDto, error) {
	if !IsPermitted(ctx, []string{"satellite.view"}) {
		//return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 3*time.Second)

	var pagination provider.Pagination[model.SatelliteLogEntry]
	dto.Map(&args, &pagination)

	pages, err := satelliteService.GetLogs(args.SatelliteId, &pagination)
	if err != nil {
		return nil, err
	}

	var pageDto LogsPageDto
	dto.Map(&pages, &pageDto)

	return &pageDto, nil
}
