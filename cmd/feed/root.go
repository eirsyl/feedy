package feed

import (
	"github.com/spf13/cobra"
)

func init() {
	FeedCmd.AddCommand(addCmd)
	FeedCmd.AddCommand(removeCmd)
	FeedCmd.AddCommand(listCmd)
}

// FeedCmd is used as the main entrypoint for the feed subcommands
var FeedCmd = &cobra.Command{
	Use:   "feed",
	Short: "Manage feeds",
	Args:  cobra.NoArgs,
}
