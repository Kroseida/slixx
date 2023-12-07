package controller

import (
	"context"
	"github.com/google/uuid"
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/reactive"
	"kroseida.org/slixx/internal/supervisor/datasource/provider"
	executionService "kroseida.org/slixx/internal/supervisor/service/execution"
	"kroseida.org/slixx/pkg/dto"
	"kroseida.org/slixx/pkg/model"
	"time"
)

type ExecutionDto struct {
	Id         uuid.UUID  `json:"id" graphql:"id"`
	JobId      uuid.UUID  `json:"jobId" graphql:"jobId"`
	Kind       string     `json:"kind" graphql:"kind"`
	Status     string     `json:"status" graphql:"status"`
	CreatedAt  time.Time  `json:"createdAt" graphql:"createdAt"`
	FinishedAt *time.Time `json:"finishedAt" graphql:"finishedAt"`
	UpdatedAt  time.Time  `json:"updatedAt" graphql:"updatedAt"`
}

type ExecutionHistoryDto struct {
	Id          uuid.UUID `json:"id" graphql:"id"`
	ExecutionId uuid.UUID `json:"executionId" graphql:"executionId"`
	Percentage  float64   `json:"percentage" graphql:"percentage"`
	StatusType  string    `json:"statusType" graphql:"statusType"`
	Message     string    `json:"message" graphql:"message"`
	CreatedAt   time.Time `json:"createdAt" graphql:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" graphql:"updatedAt"`
}

type ExecutionPageDto struct {
	Rows []ExecutionDto `json:"rows" graphql:"rows"`
	Page
}

type GetExecutionsPageDto struct {
	JobId  *uuid.UUID `json:"jobId" graphql:"jobId"`
	Limit  *int64     `json:"limit,omitempty;query:limit"`
	Page   *int64     `json:"page,omitempty;query:page"`
	Sort   *string    `json:"sort,omitempty;query:sort"`
	Search *string    `json:"search"`
}

func GetExecutions(ctx context.Context, args GetExecutionsPageDto) (*ExecutionPageDto, error) {
	if !IsPermitted(ctx, []string{"execution.view"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)

	var pagination provider.Pagination[model.Execution]
	dto.Map(&args, &pagination)

	pages, err := executionService.GetPaged(&pagination, args.JobId)
	if err != nil {
		return nil, err
	}

	var pageDto ExecutionPageDto
	dto.Map(&pages, &pageDto)

	return &pageDto, nil
}

type GetExecutionDto struct {
	ExecutionId uuid.UUID `json:"executionId" graphql:"executionId"`
}

func GetExecution(ctx context.Context, args GetExecutionDto) (*ExecutionDto, error) {
	if !IsPermitted(ctx, []string{"execution.view"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	execution, err := executionService.Get(args.ExecutionId)
	if err != nil {
		return nil, err
	}
	var executionDto *ExecutionDto
	dto.Map(&execution, &executionDto)

	return executionDto, nil
}

func GetExecutionHistory(ctx context.Context, args GetExecutionDto) ([]*ExecutionHistoryDto, error) {
	if !IsPermitted(ctx, []string{"execution.view"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	execution, err := executionService.GetHistory(args.ExecutionId)

	if err != nil {
		return nil, err
	}
	var executionDto []*ExecutionHistoryDto
	dto.Map(&execution, &executionDto)

	return executionDto, nil
}
