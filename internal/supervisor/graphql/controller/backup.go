package controller

import (
	"github.com/google/uuid"
	backupService "kroseida.org/slixx/internal/supervisor/service/backup"
)

type ExecuteBackupDto struct {
	JobId uuid.UUID `json:"jobId" graphql:"jobId"`
}

func ExecuteBackup(execute ExecuteBackupDto) (ExecuteBackupDto, error) {
	backupService.Execute(execute.JobId)
	return execute, nil
}
