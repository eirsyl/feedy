package config

import (
	"github.com/spf13/viper"
)

// GetConfigProvider returns the current config provider used for this session.
// Depends in the cli flags used when calling the program.
func GetConfigProvider() (Config, error) {
	configBackend := viper.GetString("configBackend")

	switch configBackend {
	case "file":
		return NewFileConfig()

	case "postgres":
		return NewPostgresConfig()

	default:
		return nil, ErrUnknownConfigBackend
	}
}
