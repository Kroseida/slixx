package controller

import (
	"context"
	"github.com/google/uuid"
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/reactive"
	"kroseida.org/slixx/internal/supervisor/datasource/provider"
	backupService "kroseida.org/slixx/internal/supervisor/service/backup"
	"kroseida.org/slixx/pkg/dto"
	"kroseida.org/slixx/pkg/model"
	"time"
)

type ExecuteBackupResponseDto struct {
	Id    uuid.UUID `json:"id" graphql:"id"`
	JobId uuid.UUID `json:"jobId" graphql:"jobId"`
}

type Backup struct {
	Id          uuid.UUID  `json:"id" graphql:"id"`
	Name        string     `json:"name" graphql:"name"`
	Description string     `json:"description" graphql:"description"`
	JobId       uuid.UUID  `json:"jobId" graphql:"jobId"`
	ExecutionId *uuid.UUID `json:"executionId" graphql:"executionId"`
	CreatedAt   time.Time  `json:"createdAt" graphql:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt" graphql:"updatedAt"`
}

type ExecuteBackupDto struct {
	JobId uuid.UUID `json:"jobId" graphql:"jobId"`
}

func ExecuteBackup(ctx context.Context, execute ExecuteBackupDto) (*ExecuteBackupResponseDto, error) {
	if !IsPermitted(ctx, []string{"backup.execute"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	var id, err = backupService.Execute(execute.JobId)
	if err != nil {
		return nil, err
	}
	return &ExecuteBackupResponseDto{
		Id:    *id,
		JobId: execute.JobId,
	}, nil
}

type RequestBackupSyncDto struct {
	Id uuid.UUID `json:"id" graphql:"id"`
}

func ResyncSatellite(ctx context.Context, request RequestBackupSyncDto) (*RequestBackupSyncDto, error) {
	if !IsPermitted(ctx, []string{"satellite.resync"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	var err = backupService.ResyncSatellite(request.Id)
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func DeleteBackup() {

}

type BackupPage struct {
	Rows []Backup `json:"rows" graphql:"rows"`
	Page
}

type GetBackupsRequest struct {
	JobId  *uuid.UUID `json:"jobId" graphql:"jobId"`
	Limit  *int64     `json:"limit,omitempty;query:limit"`
	Page   *int64     `json:"page,omitempty;query:page"`
	Sort   *string    `json:"sort,omitempty;query:sort"`
	Search *string    `json:"search"`
}

func GetBackups(ctx context.Context, args GetBackupsRequest) (*BackupPage, error) {
	if !IsPermitted(ctx, []string{"backup.view"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)

	var pagination provider.Pagination[model.Backup]
	dto.Map(&args, &pagination)

	pages, err := backupService.GetPaged(&pagination, args.JobId)
	if err != nil {
		return nil, err
	}

	var pageDto BackupPage
	dto.Map(&pages, &pageDto)

	return &pageDto, nil
}
