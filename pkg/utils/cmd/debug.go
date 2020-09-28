package cmd

import (
	"github.com/spf13/cobra"
)

// DebugModeEnabled returns true if the program was called with the debug flag.
func DebugModeEnabled(c *cobra.Command) bool {
	return GetBool(c, "debug")
}
