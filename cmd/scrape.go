package cmd

import (
	"github.com/spf13/viper"

	"github.com/oklog/run"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/eirsyl/feedy/pkg/client"
	"github.com/eirsyl/feedy/pkg/config"
	"github.com/eirsyl/feedy/pkg/pocket"
	"github.com/eirsyl/feedy/pkg/scraper"
	"github.com/eirsyl/feedy/pkg/worker"
	"github.com/eirsyl/flexit/cmd"
	"github.com/eirsyl/flexit/log"
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

		logger := log.NewLogrusLogger(false)

		c, err := config.NewFileConfig()
		if err != nil {
			return errors.Wrap(err, "could not create config backend")
		}
		defer c.Close() // nolint: errcheck

		cc, err := client.New()
		if err != nil {
			return errors.Wrap(err, "could not create http client")
		}

		user, err := c.GetUser()
		if err != nil {
			return errors.Wrap(err, "could not find user")
		}

		pocket, err := pocket.New(cc)
		if err != nil {
			return errors.Wrap(err, "could not create pocket client")
		}

		scr, err := scraper.New(cc)
		if err != nil {
			return errors.Wrap(err, "could not create scraper")
		}

		feeds, err := c.GetFeeds()
		if err != nil {
			return errors.Wrap(err, "could not lookup feeds")
		}

		w, err := worker.New(
			viper.GetInt("concurrency"),
			true,
			user,
			c,
			scr,
			pocket,
			logger,
		)
		if err != nil {
			return errors.Wrap(err, "could not create scrape worker")
		}

		// Add feeds to the worker queue
		for _, feed := range feeds {
			err = w.Add(feed)
			if err != nil {
				return errors.Wrap(err, "could not create feed scrape task")
			}
		}

		// Run the worker
		var g run.Group
		{
			g.Add(w.Run, w.Stop)
		}
		{
			g.Add(cmd.Interrupt(logger))
		}

		return g.Run()
	},
}
