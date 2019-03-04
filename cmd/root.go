package cmd

import (
	"github.com/eirsyl/flexit/cmd"
	"github.com/spf13/cobra"

	"github.com/eirsyl/feedy/cmd/feed"
	"github.com/eirsyl/feedy/pkg"
)

func init() {
	RootCmd.AddCommand(scrapeCmd)
	RootCmd.AddCommand(loginCmd)
	RootCmd.AddCommand(feed.FeedCmd)

	cmd.StringConfig(RootCmd, "configFile", "c", "", "config file path")
	cmd.IntConfig(RootCmd, "concurrency", "", 10, "feeds to scrape concurrent")
}

// RootCmd is ised as the main entrypoint for this application
var RootCmd = &cobra.Command{
	Use:   pkg.App.GetShortName(),
	Short: pkg.App.GetDescription(),
	Args:  cobra.NoArgs,
}
