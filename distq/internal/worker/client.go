package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	serverURL string
	http      *http.Client
}

func NewClient(serverURL string) *Client {
	return &Client{
		serverURL: strings.TrimRight(serverURL, "/"),
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type Job struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`
	Payload     json.RawMessage `json:"payload"`
	Status      string          `json:"status"`
	Attempts    int             `json:"attempts"`
	MaxAttempts int             `json:"max_attempts"`
	LeaseID     string          `json:"lease_id"`
	LeasedBy    *string         `json:"leased_by,omitempty"`
	LeaseUntil  *time.Time      `json:"lease_until,omitempty"`
}

func (c *Client) ClaimJob(ctx context.Context, workerID string, leaseSeconds int) (*Job, error) {
	reqBody := map[string]any{
		"worker_id":     workerID,
		"lease_seconds": leaseSeconds,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.serverURL+"/jobs/claim",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("claim failed: status=%s body=%s", resp.Status, string(respBody))
	}

	var job Job
	if err := json.Unmarshal(respBody, &job); err != nil {
		return nil, err
	}

	if job.LeaseID == "" {
		return nil, fmt.Errorf("claimed job has empty lease_id")
	}

	return &job, nil
}

func (c *Client) AckJob(ctx context.Context, jobID, leaseID string) error {
	body, err := json.Marshal(map[string]string{
		"lease_id": leaseID,
	})
	if err != nil {
		return err
	}

	return c.postJSON(ctx, "/jobs/"+jobID+"/ack", body)
}

func (c *Client) FailJob(ctx context.Context, jobID, leaseID, errorMsg string) error {
	body, err := json.Marshal(map[string]any{
		"lease_id": leaseID,
		"error":    errorMsg,
	})
	if err != nil {
		return err
	}

	return c.postJSON(ctx, "/jobs/"+jobID+"/fail", body)
}
