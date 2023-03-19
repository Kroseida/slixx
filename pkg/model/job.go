package model

import (
	"github.com/google/uuid"
	"time"
)

type Job struct {
	Id                   uuid.UUID `sql:"default:uuid_generate_v4()"`
	Name                 string
	Description          string
	Strategy             string
	Configuration        string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	OriginStorageId      uuid.UUID
	DestinationStorageId uuid.UUID
}
