package storage

import (
	"encoding/json"
	"time"
)

type JobStatus string

const (
	JobQueued    JobStatus = "queued"
	JobRunning   JobStatus = "running"
	JobSuccessed JobStatus = "succeeded"
	JobFailed    JobStatus = "failed"
)

type Job struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`
	Payload     json.RawMessage `json:"payload"`
	Status      JobStatus       `json:"status"`
	Attempts    int             `json:"attempts"`
	MaxAttempts int             `json:"max_attempts"`

	LeaseID    *string    `json:"lease_id,omitempty"`
	LeasedBy   *string    `json:"leased_by,omitempty"`
	LeaseUntil *time.Time `json:"lease_until,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	LastError *string `json:"last_error,omitempty"`
}
