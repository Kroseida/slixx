package model

import (
	"github.com/google/uuid"
	"time"
)

type Satellite struct {
	Id          uuid.UUID `sql:"default:uuid_generate_v4()"`
	Name        string
	Description string
	Address     string
	Token       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SatelliteLogs struct {
	Id          uuid.UUID `sql:"default:uuid_generate_v4()"`
	SatelliteId uuid.UUID
	Level       string
	Message     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}