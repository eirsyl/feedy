package feed

import (
	"net/url"

	"github.com/eirsyl/flexit/log"

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

		logger := log.NewLogrusLogger(false)

		feedURL := args[0]

		c, err := config.NewFileConfig()
		if err != nil {
			return err
		}
		defer c.Close()

		u, err := url.Parse(feedURL)
		if err != nil {
			return errors.Wrap(err, "could not parse feed url")
		}

		if err = c.RemoveFeed(u.String()); err != nil {
			return errors.Wrap(err, "could not remove feed from config")
		}

		logger.Infof("Removed feed %s from the list of feeds to scrape", u.String())

		return nil
	},
}
