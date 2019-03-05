package worker

import (
	"errors"
	"time"
)

var (
	// JobTimeout defines the maximum runtime for a task
	JobTimeout = 10 * time.Minute

	// ErrWorkerClosed Error
	ErrWorkerClosed = errors.New("Worker closed, cannot accept new jobs")
)
