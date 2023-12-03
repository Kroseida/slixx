package model

import (
	"github.com/google/uuid"
	"time"
)

type Execution struct {
	Id         uuid.UUID `sql:"default:uuid_generate_v4()"`
	Kind       string
	JobId      uuid.UUID
	Status     string
	CreatedAt  time.Time
	FinishedAt *time.Time
	UpdatedAt  time.Time
}

type ExecutionHistory struct {
	Id          uuid.UUID `sql:"default:uuid_generate_v4()"`
	ExecutionId uuid.UUID
	Percentage  float64
	StatusType  string
	Message     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
