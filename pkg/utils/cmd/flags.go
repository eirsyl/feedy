package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func flagset(cmd *cobra.Command, persistent bool) *pflag.FlagSet {
	var flags *pflag.FlagSet

	if persistent {
		flags = cmd.PersistentFlags()
	} else {
		flags = cmd.Flags()
	}

	return flags
}

// StringConfig adds a string flag to a cli.
func StringConfig(cmd *cobra.Command, name, short, value, description string, persistent bool) {
	flags := flagset(cmd, persistent)
	flags.StringP(name, short, value, description)
}

// IntConfig adds a string flag to a cli.
func IntConfig(cmd *cobra.Command, name, short string, value int, description string, persistent bool) {
	flags := flagset(cmd, persistent)
	flags.IntP(name, short, value, description)
}

// BoolConfig adds a bool flag to a cli.
func BoolConfig(cmd *cobra.Command, name, short string, value bool, description string, persistent bool) {
	flags := flagset(cmd, persistent)
	flags.BoolP(name, short, value, description)
}

// DurationConfig adds a duration flag to a cli.
func DurationConfig(cmd *cobra.Command, name, short string, value time.Duration, description string, persistent bool) {
	flags := flagset(cmd, persistent)
	flags.DurationP(name, short, value, description)
}

// StringSliceConfig adds a string slice flag to a cli.
func StringSliceConfig(cmd *cobra.Command, name, short string, value []string, description string, persistent bool) {
	flags := flagset(cmd, persistent)
	flags.StringSliceP(name, short, value, description)
}

// StringToStringConfig adds a string to string (string map) flag to a cli.
func StringToStringConfig(
	cmd *cobra.Command, name, short string, value map[string]string,
	description string, persistent bool,
) {
	flags := flagset(cmd, persistent)
	flags.StringToStringP(name, short, value, description)
}
