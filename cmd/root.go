package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/eirsyl/feedy/cmd/feed"
	"github.com/eirsyl/feedy/internal"
	"github.com/eirsyl/feedy/pkg/utils/cmd"
	"github.com/eirsyl/feedy/pkg/utils/runtime"
)

func init() {
	cmd.StringConfig(RootCmd, "log-level", "", "info", "log level", true)
	cmd.StringConfig(RootCmd, "log-format", "", "text", "log format", true)
	cmd.BoolConfig(RootCmd, "debug", "", false, "run program in debug mode", true)

	RootCmd.AddCommand(scrapeCmd)
	RootCmd.AddCommand(loginCmd)
	RootCmd.AddCommand(feed.FeedCmd)
}

// RootCmd is ised as the main entrypoint for this application
var RootCmd = &cobra.Command{
	Use:     "feedy",
	Short:   internal.ShortDescription,
	Long:    internal.LongDescription,
	Args:    cobra.NoArgs,
	Version: fmt.Sprintf("%s (%s)", internal.Version, internal.BuildDate),
	PersistentPreRunE: func(c *cobra.Command, _ []string) error {
		// Initialize logger based on command line flags
		logger, err := cmd.InitializeLogger(c)
		if err != nil {
			return err
		}

		return runtime.OptimizeRuntime(logger)
	},
	RunE: func(c *cobra.Command, args []string) error {
		return c.Help()
	},
}
