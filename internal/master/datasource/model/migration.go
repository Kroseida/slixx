package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Migration struct {
	gorm.Model
	Id   uuid.UUID `sql:"default:uuid_generate_v4()"`
	Name string
}
