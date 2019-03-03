package feed

import (
	"github.com/eirsyl/flexit/log"

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

		logger := log.NewLogrusLogger(false)

		c, err := config.NewFileConfig()
		if err != nil {
			return err
		}
		defer c.Close()

		feeds, err := c.GetFeeds()
		if err != nil {
			return errors.Wrap(err, "could not retrieve feeds")
		}

		for _, f := range feeds {
			logger.Infof("%s %s %s", f.Name, f.URL, f.Tags)
		}

		return nil
	},
}
