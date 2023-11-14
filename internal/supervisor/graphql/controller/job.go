package controller

import (
	"context"
	"github.com/google/uuid"
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/reactive"
	"kroseida.org/slixx/internal/supervisor/datasource/provider"
	jobService "kroseida.org/slixx/internal/supervisor/service/job"
	"kroseida.org/slixx/pkg/dto"
	"kroseida.org/slixx/pkg/model"
	"kroseida.org/slixx/pkg/strategy"
	"reflect"
	"time"
)

type Job struct {
	Id                   uuid.UUID `sql:"default:uuid_generate_v4()"`
	Name                 string
	Description          string
	Strategy             string
	Configuration        string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	OriginStorageId      uuid.UUID
	DestinationStorageId uuid.UUID
	ExecutorSatelliteId  uuid.UUID
}

type GetJobDto struct {
	Id uuid.UUID `json:"id" graphql:"id"`
}

type JobsPage struct {
	Rows []Job `json:"rows" graphql:"rows"`
	Page
}

func GetJob(ctx context.Context, args GetJobDto) (*Job, error) {
	if !IsPermitted(ctx, []string{"job.view"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	job, err := jobService.Get(args.Id)
	if err != nil {
		return nil, err
	}
	var jobDto *Job
	dto.Map(&job, &jobDto)

	return jobDto, nil
}

func GetJobs(ctx context.Context, args PageArgs) (*JobsPage, error) {
	if !IsPermitted(ctx, []string{"job.view"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)

	var pagination provider.Pagination[model.Job]
	dto.Map(&args, &pagination)

	pages, err := jobService.GetPaged(&pagination)
	if err != nil {
		return nil, err
	}

	var pageDto JobsPage
	dto.Map(&pages, &pageDto)

	return &pageDto, nil
}

type CreateJobDto struct {
	Name                 string
	Description          string
	Strategy             string
	Configuration        string
	OriginStorageId      uuid.UUID
	DestinationStorageId uuid.UUID
	ExecutorSatelliteId  uuid.UUID
}

func CreateJob(ctx context.Context, args CreateJobDto) (*Job, error) {
	if !IsPermitted(ctx, []string{"job.create"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	job, err := jobService.Create(
		args.Name,
		args.Description,
		args.Strategy,
		args.Configuration,
		args.OriginStorageId,
		args.DestinationStorageId,
		args.ExecutorSatelliteId,
	)
	if err != nil {
		return nil, err
	}
	var jobsDto Job
	dto.Map(job, &jobsDto)

	return &jobsDto, err
}

type UpdateJobDto struct {
	Id                   uuid.UUID
	Name                 *string
	Description          *string
	Strategy             *string
	Configuration        *string
	OriginStorageId      *uuid.UUID
	DestinationStorageId *uuid.UUID
	ExecutorSatelliteId  *uuid.UUID
}

func UpdateJob(ctx context.Context, args UpdateJobDto) (*Job, error) {
	if !IsPermitted(ctx, []string{"job.update"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	job, err := jobService.Update(
		args.Id,
		args.Name,
		args.Description,
		args.Strategy,
		args.Configuration,
		args.OriginStorageId,
		args.DestinationStorageId,
		args.ExecutorSatelliteId,
	)
	if err != nil {
		return nil, err
	}
	var jobDto Job
	dto.Map(job, &jobDto)

	return &jobDto, err
}

type DeleteJobDto struct {
	Id uuid.UUID `json:"id" graphql:"id"`
}

func DeleteJob(ctx context.Context, args DeleteJobDto) (*Job, error) {
	if !IsPermitted(ctx, []string{"job.delete"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	job, err := jobService.Delete(args.Id)
	if err != nil {
		return nil, err
	}
	var jobDto Job
	dto.Map(job, &jobDto)

	return &jobDto, err
}

type JobStrategyDescriptionDto struct {
	Name          string                                   `json:"name" graphql:"name"`
	Configuration []JobStrategyConfigurationDescriptionDto `json:"configuration" graphql:"configuration"`
}

type JobStrategyConfigurationDescriptionDto struct {
	Name    string `json:"name" graphql:"name"`
	Kind    string `json:"kind" graphql:"kind"`
	Default string `json:"default" graphql:"default"`
}

func GetJobStrategies() ([]JobStrategyDescriptionDto, error) {
	var descriptions []JobStrategyDescriptionDto
	for _, kind := range strategy.Values() {
		var configurations []JobStrategyConfigurationDescriptionDto

		val := reflect.ValueOf(kind.DefaultConfiguration()).Elem()
		for i := 0; i < val.NumField(); i++ {
			configurations = append(configurations, JobStrategyConfigurationDescriptionDto{
				Name:    val.Type().Field(i).Tag.Get("json"),
				Kind:    val.Type().Field(i).Tag.Get("slixx"),
				Default: val.Type().Field(i).Tag.Get("default"),
			})
		}

		descriptions = append(descriptions, JobStrategyDescriptionDto{
			Name:          kind.GetName(),
			Configuration: configurations,
		})
	}

	return descriptions, nil
}
