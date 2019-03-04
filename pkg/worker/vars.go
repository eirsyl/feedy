package worker

import "errors"

var (
	// ErrWorkerClosed Error
	ErrWorkerClosed = errors.New("Worker closed, cannot accept new jobs")
)
