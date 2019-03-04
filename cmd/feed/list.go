package feed

import (
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/eirsyl/feedy/pkg/config"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list feeds",
	Long: `
List all feeds configured for scraping.
	`,
	Args: cobra.NoArgs,
	PreRun: func(_ *cobra.Command, args []string) {
	},
	RunE: func(_ *cobra.Command, args []string) error {

		c, err := config.NewFileConfig()
		if err != nil {
			return err
		}
		defer c.Close() // nolint: errcheck

		feeds, err := c.GetFeeds()
		if err != nil {
			return errors.Wrap(err, "could not retrieve feeds")
		}

		tableWriter := table.NewWriter()
		tableWriter.SetOutputMirror(os.Stdout)
		tableWriter.SetStyle(table.StyleRounded)
		tableWriter.AppendHeader(table.Row{"Name", "URL", "Tags"})

		for _, f := range feeds {
			tableWriter.AppendRow(table.Row{f.Name, f.URL, strings.Join(f.Tags, ", ")})
		}

		tableWriter.AppendFooter(table.Row{"", "Total", len(feeds)})
		tableWriter.Render()

		return nil
	},
}
