package parallel

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/pkg/utils"
	"time"
)

type RunningJob struct {
	JobId     uuid.UUID
	Canceled  bool
	Callback  func(update utils.StatusUpdate)
	StartedAt time.Time
}

func (job *RunningJob) CheckCancel() bool {
	return job.Canceled
}
