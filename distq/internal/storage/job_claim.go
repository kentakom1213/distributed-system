package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

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
