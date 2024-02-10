package controller

import (
	"context"
	"github.com/google/uuid"
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/reactive"
	"kroseida.org/slixx/internal/supervisor/datasource/provider"
	jobScheduleService "kroseida.org/slixx/internal/supervisor/service/job_schedule"
	"kroseida.org/slixx/pkg/dto"
	"kroseida.org/slixx/pkg/model"
	"kroseida.org/slixx/pkg/schedule"
	"reflect"
	"time"
)

type JobSchedule struct {
	Id            uuid.UUID `json:"id" graphql:"id"`
	Name          string    `json:"name" graphql:"name"`
	Description   string    `json:"description" graphql:"description"`
	Kind          string    `json:"kind" graphql:"kind"`
	Configuration string    `json:"configuration" graphql:"configuration"`
	CreatedAt     time.Time `json:"createdAt" graphql:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt" graphql:"updatedAt"`
	DeletedAt     time.Time `json:"deletedAt" graphql:"deletedAt"`
}

type JobSchedulePrototypeDto struct {
	Id          uuid.UUID `json:"id" graphql:"id"`
	Name        string    `json:"name" graphql:"name"`
	Description string    `json:"description" graphql:"description"`
	Kind        string    `json:"kind" graphql:"kind"`
	CreatedAt   time.Time `json:"createdAt" graphql:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" graphql:"updatedAt"`
	DeletedAt   time.Time `json:"deletedAt" graphql:"deletedAt"`
}

type JobSchedulesPageDto struct {
	Rows []JobSchedulePrototypeDto `json:"rows" graphql:"rows"`
	Page
}

type GetJobScheduleDto struct {
	Id uuid.UUID `json:"id" graphql:"id"`
}

func GetJobSchedule(ctx context.Context, args GetJobScheduleDto) (*JobSchedule, error) {
	if !IsPermitted(ctx, []string{"jobSchedule.view"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	jobSchedule, err := jobScheduleService.Get(args.Id)
	if err != nil {
		return nil, err
	}
	var jobScheduleDto *JobSchedule
	dto.Map(&jobSchedule, &jobScheduleDto)

	return jobScheduleDto, nil
}

type GetJobSchedulesPageRequest struct {
	JobId  *uuid.UUID `json:"jobId" graphql:"jobId"`
	Limit  *int64     `json:"limit,omitempty;query:limit"`
	Page   *int64     `json:"page,omitempty;query:page"`
	Sort   *string    `json:"sort,omitempty;query:sort"`
	Search *string    `json:"search"`
}

func GetJobSchedules(ctx context.Context, args GetJobSchedulesPageRequest) (*JobSchedulesPageDto, error) {
	if !IsPermitted(ctx, []string{"jobSchedule.view"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)

	var pagination provider.Pagination[model.JobSchedule]
	dto.Map(&args, &pagination)

	pages, err := jobScheduleService.GetPaged(&pagination, args.JobId)
	if err != nil {
		return nil, err
	}

	var pageDto JobSchedulesPageDto
	dto.Map(&pages, &pageDto)

	return &pageDto, nil
}

type CreateJobScheduleDto struct {
	Name          string `json:"name" graphql:"name"`
	Description   string `json:"description" graphql:"description"`
	Kind          string `json:"kind" graphql:"kind"`
	Configuration string `json:"configuration" graphql:"configuration"`
}

func CreateJobSchedule(ctx context.Context, args CreateJobScheduleDto) (*JobSchedule, error) {
	if !IsPermitted(ctx, []string{"jobSchedule.create"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	jobSchedule, err := jobScheduleService.Create(args.Name, args.Description, args.Kind, args.Configuration)
	if err != nil {
		return nil, err
	}
	var jobScheduleDto JobSchedule
	dto.Map(jobSchedule, &jobScheduleDto)

	return &jobScheduleDto, err
}

type UpdateJobScheduleDto struct {
	Id            uuid.UUID `json:"id" graphql:"id"`
	Name          *string   `json:"name" graphql:"name"`
	Description   *string   `json:"description" graphql:"description"`
	Kind          *string   `json:"kind" graphql:"kind"`
	Configuration *string   `json:"configuration" graphql:"configuration"`
}

func UpdateJobSchedule(ctx context.Context, args UpdateJobScheduleDto) (*JobSchedule, error) {
	if !IsPermitted(ctx, []string{"jobSchedule.update"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	jobScheduleModel, err := jobScheduleService.Update(
		args.Id,
		args.Name,
		args.Description,
		args.Kind,
		args.Configuration,
	)

	if err != nil {
		return nil, err
	}
	var jobScheduleDto JobSchedule
	dto.Map(&jobScheduleModel, &jobScheduleDto)

	return &jobScheduleDto, nil
}

type DeleteJobScheduleDto struct {
	Id uuid.UUID `json:"id" graphql:"id"`
}

func DeleteJobSchedule(ctx context.Context, args DeleteJobScheduleDto) (*JobSchedule, error) {
	if !IsPermitted(ctx, []string{"jobSchedule.delete"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	jobSchedule, err := jobScheduleService.Delete(args.Id)
	if err != nil {
		return nil, err
	}
	var jobScheduleDto JobSchedule
	dto.Map(&jobSchedule, &jobScheduleDto)

	return &jobScheduleDto, nil
}

type JobScheduleKindDescriptionDto struct {
	Name          string                                   `json:"name" graphql:"name"`
	Configuration []JobScheduleConfigurationDescriptionDto `json:"configuration" graphql:"configuration"`
}

type JobScheduleConfigurationDescriptionDto struct {
	Name    string `json:"name" graphql:"name"`
	Kind    string `json:"kind" graphql:"kind"`
	Default string `json:"default" graphql:"default"`
}

func GetJobScheduleKinds() ([]JobScheduleKindDescriptionDto, error) {
	var descriptions []JobScheduleKindDescriptionDto
	for _, kind := range schedule.Values() {
		var configurations []JobScheduleConfigurationDescriptionDto

		val := reflect.ValueOf(kind.DefaultConfiguration()).Elem()
		for i := 0; i < val.NumField(); i++ {
			configurations = append(configurations, JobScheduleConfigurationDescriptionDto{
				Name:    val.Type().Field(i).Tag.Get("json"),
				Kind:    val.Type().Field(i).Tag.Get("slixx"),
				Default: val.Type().Field(i).Tag.Get("default"),
			})
		}

		descriptions = append(descriptions, JobScheduleKindDescriptionDto{
			Name:          kind.GetName(),
			Configuration: configurations,
		})
	}
	return descriptions, nil
}
