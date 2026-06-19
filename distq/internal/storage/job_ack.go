package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (s *Storage) AckJob(ctx context.Context, jobID, leaseID string) error {
	if jobID == "" {
		return fmt.Errorf("job id is required")
	}
	if leaseID == "" {
		return fmt.Errorf("lease id is required")
	}

	now := nowJST()
	nowText := formatTime(now)

	res, err := s.db.ExecContext(ctx, `
		UPDATE jobs
		SET
			status = ?,
			lease_id = NULL,
			leased_by = NULL,
			lease_until = NULL,
			updated_at = ?
		WHERE
			id = ?
			AND lease_id = ?
			AND status = ?
			AND lease_until IS NOT NULL
			AND lease_until >= ?
	`,
		string(JobSuccessed),
		nowText,
		jobID,
		leaseID,
		string(JobRunning),
		nowText,
	)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 1 {
		return nil
	}

	return s.explainAckFailure(ctx, jobID, leaseID, now)
}

func (s *Storage) explainAckFailure(ctx context.Context, jobID, leaseID string, now time.Time) error {
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
		return ErrJobNotFound
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
