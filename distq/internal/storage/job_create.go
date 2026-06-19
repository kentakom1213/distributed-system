package storage

import (
	"context"
	"encoding/json"
	"fmt"
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
