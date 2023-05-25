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
	ExecutorSatelliteId  uuid.UUID
	// We have to check that (if origin is "file system" storage type) the origin storage is on the same satellite as the job
	// So we can also create an archive on origin satellite and then transfer it to destination satellite
	// while crating partial backup we send the hash of the file to the destination satellite and then the destination satellite checks if it has the file
	// if it has the file it will not download it again from the origin satellite
}
