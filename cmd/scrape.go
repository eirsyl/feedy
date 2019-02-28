package cmd

import (
	"github.com/spf13/cobra"
)

var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Scrape watched feeds",
	Long: `
Scrape watched feeds and send new items to pocket.
	`,
	Args: cobra.NoArgs,
	PreRun: func(_ *cobra.Command, args []string) {
	},
	RunE: func(_ *cobra.Command, args []string) error {
		return nil
	},
}
