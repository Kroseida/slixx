package utils

import "github.com/google/uuid"

type StatusUpdate struct {
	Id         uuid.UUID  `json:"id"`
	JobId      *uuid.UUID `json:"jobId"`
	Percentage float64    `json:"percentage"`
	Message    string     `json:"message"`
	StatusType string     `json:"statusType"`
}
