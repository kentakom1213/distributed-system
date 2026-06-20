package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (s *Storage) FailJob(
	ctx context.Context,
	jobID, leaseID, errorMsg string,
) (*Job, error) {
	if jobID == "" {
		return nil, fmt.Errorf("job id is required")
	}
	if leaseID == "" {
		return nil, fmt.Errorf("lease id is required")
	}
	if errorMsg == "" {
		errorMsg = "job failed"
	}

	now := nowJST()
	nowText := formatTime(now)

	row := s.db.QueryRowContext(ctx, `
		UPDATE jobs
		SET
			status = CASE
				WHEN attempts >= max_attempts THEN ?
				ELSE ?
			END,
			lease_id = NULL,
			leased_by = NULL,
			lease_until = NULL,
			updated_at = ?,
			last_error = ?
		WHERE
			id = ?
			AND lease_id = ?
			AND status = ?
			AND lease_until IS NOT NULL
			AND lease_until >= ?
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
		string(JobFailed),
		string(JobQueued),
		nowText,
		errorMsg,
		jobID,
		leaseID,
		string(JobRunning),
		nowText,
	)

	job, err := scanJob(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, s.explainFailFailure(ctx, jobID, leaseID, now)
	}
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (s *Storage) explainFailFailure(ctx context.Context, jobID, leaseID string, now time.Time) error {
	var status string
	var currentLeaseID sql.NullString
	var leaseUntil sql.NullString

	err := s.db.QueryRowContext(ctx, `
		SELECT
			status,
			lease_id,
			lease_until
		FROM jobs
		WHERE id = ?
	`, jobID).Scan(
		&status,
		&currentLeaseID,
		&leaseUntil,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return ErrJobNotFound
	}
	if err != nil {
		return err
	}

	if JobStatus(status) != JobRunning {
		return ErrJobNotRunning
	}

	if !currentLeaseID.Valid || currentLeaseID.String != leaseID {
		return ErrInvalidLease
	}

	if !leaseUntil.Valid {
		return ErrInvalidLease
	}

	parsedLeaseUntil, err := parseTime(leaseUntil.String)
	if err != nil {
		return err
	}

	if parsedLeaseUntil.Before(now) {
		return ErrLeaseExpired
	}

	return ErrInvalidLease
}
