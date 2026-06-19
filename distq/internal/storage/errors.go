package storage

import "errors"

var (
	ErrJobNotFound   = errors.New("job not found")
	ErrInvalidLease  = errors.New("invalid lease")
	ErrLeaseExpired  = errors.New("lease expired")
	ErrJobNotRunning = errors.New("job is not running")
)
