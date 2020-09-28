package cmd

import (
	"errors"

	"github.com/eirsyl/feedy/pkg/utils/log"
	"github.com/spf13/cobra"
)

var (
	// ErrInvalidLogLevel is raised the the provided log level is invalid
	ErrInvalidLogLevel = errors.New("invalid log level")
	// ErrInvalidLogFormat is raised the the provided log format is invalid
	ErrInvalidLogFormat = errors.New("invalid log format")
)

// InitializeLogger in initializes a new logger based on the command line flags
// provided when the program started.
func InitializeLogger(c *cobra.Command) (log.Logger, error) {
	var logLevel, logFormat string
	{
		logLevel = GetString(c, "log-level")
		logFormat = GetString(c, "log-format")
	}

	var logger log.Logger

	// Create logger based on log format
	switch logFormat {
	case "text":
		logger = log.NewLogrusLogger(false)
	case "json":
		logger = log.NewLogrusLogger(true)
	default:
		return nil, ErrInvalidLogFormat
	}

	// Set level debug if program was called with the debug flag
	if DebugModeEnabled(c) {
		logLevel = "debug"
	}

	// Configure log level
	err := setLogLevel(logger, logLevel)
	if err != nil {
		return nil, err
	}

	return logger, nil
}

func setLogLevel(logger log.Logger, level string) error {
	switch level {
	case "debug":
		logger.SetLevel(log.DebugLevel)
	case "info":
		logger.SetLevel(log.InfoLevel)
	case "warn":
		logger.SetLevel(log.WarnLevel)
	case "error":
		logger.SetLevel(log.ErrorLevel)
	case "fatal":
		logger.SetLevel(log.FatalLevel)
	case "panic":
		logger.SetLevel(log.PanicLevel)
	default:
		return ErrInvalidLogLevel
	}

	return nil
}
