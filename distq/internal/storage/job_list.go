package storage

import (
	"context"
	"database/sql"
)

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
