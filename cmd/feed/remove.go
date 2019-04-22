package feed

import (
	"fmt"
	"net/url"

	"github.com/jedib0t/go-pretty/text"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/eirsyl/feedy/pkg/config"
)

var removeCmd = &cobra.Command{
	Use:   "remove [feed]",
	Short: "Remove feed",
	Long: `
Remove feed from the list of feeds to scrape.
	`,
	Args: cobra.ExactArgs(1),
	PreRun: func(_ *cobra.Command, args []string) {
	},
	RunE: func(_ *cobra.Command, args []string) error {

		feedURL := args[0]

		c, err := config.GetConfigProvider()
		if err != nil {
			return err
		}
		defer c.Close() // nolint: errcheck

		u, err := url.Parse(feedURL)
		if err != nil {
			return errors.Wrap(err, "could not parse feed url")
		}

		if err = c.RemoveFeed(u.String()); err != nil {
			return errors.Wrap(err, "could not remove feed from config")
		}

		fmt.Println(text.FgGreen.Sprintf("Successfully removed %s from the watch list", u.String()))

		return nil
	},
}
