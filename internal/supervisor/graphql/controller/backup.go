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

type ExecuteBackupDto struct {
	JobId uuid.UUID `json:"jobId" graphql:"jobId"`
}

type Backup struct {
	Id          uuid.UUID `sql:"default:uuid_generate_v4()"`
	Name        string
	Description string
	JobId       uuid.UUID
	ExecutionId *uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func ExecuteBackup(execute ExecuteBackupDto) (ExecuteBackupDto, error) {
	backupService.Execute(execute.JobId)
	return execute, nil
}

type BackupPage struct {
	Rows []Backup `json:"rows" graphql:"rows"`
	Page
}

func GetBackups(ctx context.Context, args PageArgs) (*BackupPage, error) {
	if !IsPermitted(ctx, []string{"backup.view"}) {
		return nil, graphql.NewSafeError("missing permission")
	}
	reactive.InvalidateAfter(ctx, 5*time.Second)

	var pagination provider.Pagination[model.Backup]
	dto.Map(&args, &pagination)

	pages, err := backupService.GetPaged(&pagination)
	if err != nil {
		return nil, err
	}

	var pageDto BackupPage
	dto.Map(&pages, &pageDto)

	return &pageDto, nil
}
