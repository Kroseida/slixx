package model

import (
	"github.com/google/uuid"
	"time"
)

// Backup
// The Backup does contain the information about his job. In this way we can recreate the job from the backup.
// e.g if the job is deleted or corrupted, we can recreate it from the backup.
type Backup struct {
	Id              uuid.UUID `sql:"default:uuid_generate_v4()"`
	Name            string
	Description     string
	JobId           uuid.UUID
	ExecutionId     uuid.UUID
	OriginKind      string
	DestinationKind string
	Strategy        string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
