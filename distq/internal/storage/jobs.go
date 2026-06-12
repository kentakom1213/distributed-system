package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type CreateJobParams struct {
	Type        string
	Payload     json.RawMessage
	MaxAttempts int
}

func (s *Storage) CreateJob(ctx context.Context, p CreateJobParams) (*Job, error) {
	if p.Type == "" {
		return nil, fmt.Errorf("job type is required")
	}

	if !json.Valid(p.Payload) {
		return nil, fmt.Errorf("payload must be valid JSON")
	}

	if p.MaxAttempts <= 0 {
		p.MaxAttempts = 3
	}

	id, err := randomID("job")
	if err != nil {
		return nil, err
	}

	now := nowJST()
	nowText := formatTime(now)

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO jobs (
			id,
			type,
			payload,
			status,
			attempts,
			max_attempts,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		id,
		p.Type,
		string(p.Payload),
		string(JobQueued),
		0,
		p.MaxAttempts,
		nowText,
		nowText,
	)
	if err != nil {
		return nil, err
	}

	return &Job{
		ID:          id,
		Type:        p.Type,
		Payload:     p.Payload,
		Status:      JobQueued,
		Attempts:    0,
		MaxAttempts: p.MaxAttempts,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

type ListJobsParams struct {
	Status JobStatus
	Limit  int
}

func (s *Storage) ListJobs(ctx context.Context, p ListJobsParams) ([]Job, error) {
	if p.Limit <= 0 {
		p.Limit = 100
	}

	var rows *sql.Rows
	var err error

	if p.Status != "" {
		rows, err = s.db.QueryContext(ctx, `
			SELECT
				id,
				type,
				payload,
				status,
				attempts,
				max_attempts,
				lease_id,
				leased_by,
				lease_until,
				created_at,
				updated_at,
				last_error
			FROM jobs
			WHERE status = ?
			ORDER BY created_at ASC
			LIMIT ?
		`, string(p.Status), p.Limit)
	} else {
		rows, err = s.db.QueryContext(ctx, `
			SELECT
				id,
				type,
				payload,
				status,
				attempts,
				max_attempts,
				lease_id,
				leased_by,
				lease_until,
				created_at,
				updated_at,
				last_error
			FROM jobs
			ORDER BY created_at ASC
			LIMIT ?
		`, p.Limit)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []Job

	for rows.Next() {
		job, err := scanJob(rows)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, job)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

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

type PickNextJobParams struct {
	WorkerID      string
	LeaseDuration time.Duration
}

func (s *Storage) PickNextJob(ctx context.Context, p PickNextJobParams) (*Job, error) {
	if p.WorkerID == "" {
		return nil, fmt.Errorf("worker_id is required")
	}

	if p.LeaseDuration <= 0 {
		p.LeaseDuration = 30 * time.Second
	}

	leaseID, err := randomID("lease")
	if err != nil {
		return nil, err
	}

	now := nowJST()
	leaseUntil := now.Add(p.LeaseDuration)

	nowText := formatTime(now)
	leaseUntilText := formatTime(leaseUntil)

	row := s.db.QueryRowContext(ctx, `
		UPDATE jobs
		SET
			status = ?,
			attempts = attempts + 1,
			lease_id = ?,
			leased_by = ?,
			lease_until = ?,
			updated_at = ?
		WHERE id = (
			SELECT id
			FROM jobs
			WHERE
				(
					status = ?
					AND attempts < max_attempts
				)
				OR (
					status = ?
					AND lease_until IS NOT NULL
					AND lease_until < ?
					AND attempts < max_attempts
				)
			ORDER BY created_at ASC
			LIMIT 1
		)
		RETURNING
			id,
			type,
			payload,
			status,
			attempts,
			max_attempts,
			lease_id,
			leased_by,
			lease_until,
			created_at,
			updated_at,
			last_error
	`,
		string(JobRunning),
		leaseID,
		p.WorkerID,
		leaseUntilText,
		nowText,
		string(JobQueued),
		string(JobRunning),
		nowText,
	)

	job, err := scanJob(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &job, nil
}
