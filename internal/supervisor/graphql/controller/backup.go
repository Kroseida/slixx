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

type BackupDto struct {
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

type RestoreBackupDto struct {
	BackupId uuid.UUID `json:"backupId" graphql:"backupId"`
}

type RestoreBackupResponseDto struct {
	Id uuid.UUID `json:"id" graphql:"id"`
}

func RestoreBackup(ctx context.Context, restore RestoreBackupDto) (*RestoreBackupResponseDto, error) {
	if !IsPermitted(ctx, []string{"backup.restore"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	var id, err = backupService.Restore(restore.BackupId)
	if err != nil {
		return nil, err
	}
	return &RestoreBackupResponseDto{
		Id: *id,
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

type BackupPageDto struct {
	Rows []BackupDto `json:"rows" graphql:"rows"`
	Page
}

type GetBackupsPageRequest struct {
	JobId  *uuid.UUID `json:"jobId" graphql:"jobId"`
	Limit  *int64     `json:"limit,omitempty;query:limit"`
	Page   *int64     `json:"page,omitempty;query:page"`
	Sort   *string    `json:"sort,omitempty;query:sort"`
	Search *string    `json:"search"`
}

func GetBackups(ctx context.Context, args GetBackupsPageRequest) (*BackupPageDto, error) {
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

	var pageDto BackupPageDto
	dto.Map(&pages, &pageDto)

	return &pageDto, nil
}

type GetBackupDto struct {
	Id uuid.UUID `json:"id" graphql:"id"`
}

func GetBackup(ctx context.Context, args GetBackupDto) (*BackupDto, error) {
	if !IsPermitted(ctx, []string{"backup.view"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)
	job, err := backupService.Get(args.Id)
	if err != nil {
		return nil, err
	}
	var backupDto *BackupDto
	dto.Map(&job, &backupDto)

	return backupDto, nil
}

type DeleteBackupDto struct {
	Id       *uuid.UUID `json:"id" graphql:"id"`
	JobId    uuid.UUID  `json:"jobId" graphql:"jobId"`
	BackupId uuid.UUID  `json:"backupId" graphql:"backupId"`
}

func DeleteBackup(ctx context.Context, args DeleteBackupDto) (*DeleteBackupDto, error) {
	if !IsPermitted(ctx, []string{"backup.delete"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	if args.Id == nil {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, err
		}
		args.Id = &id
	}
	var err = backupService.RequestDelete(*args.Id, args.JobId, args.BackupId)
	if err != nil {
		return nil, err
	}
	return &args, nil
}
