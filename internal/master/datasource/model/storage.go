package model

import (
	"github.com/google/uuid"
	"time"
)

type Storage struct {
	Id            uuid.UUID `sql:"default:uuid_generate_v4()"`
	Name          string
	Kind          string
	Configuration string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
