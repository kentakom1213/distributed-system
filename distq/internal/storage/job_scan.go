package storage

import (
	"database/sql"
	"encoding/json"
)

type jobScanner interface {
	Scan(dest ...any) error
}

func scanJob(scanner jobScanner) (Job, error) {
	var (
		job     Job
		payload string
		status  string

		leaseID    sql.NullString
		leasedBy   sql.NullString
		leaseUntil sql.NullString
		createdAt  string
		updatedAt  string
		lastError  sql.NullString
	)

	err := scanner.Scan(
		&job.ID,
		&job.Type,
		&payload,
		&status,
		&job.Attempts,
		&job.MaxAttempts,
		&leaseID,
		&leasedBy,
		&leaseUntil,
		&createdAt,
		&updatedAt,
		&lastError,
	)
	if err != nil {
		return Job{}, err
	}

	job.Payload = json.RawMessage(payload)
	job.Status = JobStatus(status)

	parsedCreatedAt, err := parseTime(createdAt)
	if err != nil {
		return Job{}, err
	}
	job.CreatedAt = parsedCreatedAt

	parsedUpdatedAt, err := parseTime(updatedAt)
	if err != nil {
		return Job{}, err
	}
	job.UpdatedAt = parsedUpdatedAt

	if leaseID.Valid {
		job.LeaseID = &leaseID.String
	}
	if leasedBy.Valid {
		job.LeasedBy = &leasedBy.String
	}
	if leaseUntil.Valid {
		t, err := parseTime(leaseUntil.String)
		if err != nil {
			return Job{}, err
		}
		job.LeaseUntil = &t
	}
	if lastError.Valid {
		job.LastError = &lastError.String
	}

	return job, nil
}
