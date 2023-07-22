package model

import (
	"github.com/google/uuid"
	"time"
)

type Backup struct {
	Id          uuid.UUID `sql:"default:uuid_generate_v4()"`
	Name        string
	Description string
	JobId       uuid.UUID
	ExecutionId uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
