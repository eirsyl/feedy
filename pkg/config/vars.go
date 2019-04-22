package config

import "errors"

var (
	// ErrConfigFileNotGiven Error
	ErrConfigFileNotGiven = errors.New("Config file path not provided")

	// ErrFeedNotFound Error
	ErrFeedNotFound = errors.New("Feed not found")

	// ErrUnknownConfigBackend Error
	ErrUnknownConfigBackend = errors.New("unknown config backend")
)
