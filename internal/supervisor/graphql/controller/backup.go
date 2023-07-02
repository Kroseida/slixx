package controller

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/internal/supervisor/syncnetwork/action"
)

type ExecuteBackupDto struct {
	JobId uuid.UUID `json:"jobId" graphql:"jobId"`
}

func ExecuteBackup(execute ExecuteBackupDto) (ExecuteBackupDto, error) {
	action.SendExecuteBackup(execute.JobId)
	return execute, nil
}
