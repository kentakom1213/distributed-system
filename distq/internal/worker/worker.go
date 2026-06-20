package worker

import "time"

type Config struct {
	ID           string
	ServerURL    string
	PollInterval string
	LeaseSeconds int
}

type Worker struct {
	id           string
	client       *Client
	pollInterval time.Duration
	leaseSeconds int
}
