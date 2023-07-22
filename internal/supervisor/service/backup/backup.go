package backup

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/internal/supervisor/syncnetwork/action"
)

func Execute(jobId uuid.UUID) {
	action.SendExecuteBackup(jobId)
}
