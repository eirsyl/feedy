package cmd

import (
	"time"

	"github.com/oklog/run"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/eirsyl/feedy/pkg/client"
	"github.com/eirsyl/feedy/pkg/config"
	"github.com/eirsyl/feedy/pkg/pocket"
	"github.com/eirsyl/feedy/pkg/scraper"
	"github.com/eirsyl/feedy/pkg/worker"
	"github.com/eirsyl/flexit/cmd"
	"github.com/eirsyl/flexit/log"
)

func init() {
	cmd.BoolConfig(scrapeCmd, "autostop", "", true, "stop scraping of feeds after one iteration")
}

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

		var autostop bool
		{
			autostop = viper.GetBool("autostop")
		}

		c, err := config.GetConfigProvider()
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

		w, err := worker.New(
			viper.GetInt("concurrency"),
			user,
			c,
			scr,
			pocket,
			logger,
		)
		if err != nil {
			return errors.Wrap(err, "could not create scrape worker")
		}

		s := scheduler{
			w:        w,
			c:        c,
			autostop: autostop,
			stopped:  make(chan struct{}, 1),
		}

		// Run the worker
		var g run.Group
		{
			g.Add(w.Run, w.Stop)
		}
		{
			g.Add(s.run, s.stop)
		}
		{
			g.Add(cmd.Interrupt(logger))
		}

		return g.Run()
	},
}

type scheduler struct {
	w        worker.Worker
	c        config.Config
	autostop bool
	stopped  chan struct{}
}

func (s *scheduler) run() error {
	for {

		// Lookup feeds from config
		feeds, err := s.c.GetFeeds()
		if err != nil {
			return errors.Wrap(err, "could not lookup feeds")
		}

		// Add feeds to the worker queue
		for _, feed := range feeds {
			err = s.w.Add(feed)
			if err != nil {
				return errors.Wrap(err, "could not create feed scrape task")
			}
		}

		if s.autostop {
			return nil
		}

		select {
		case <-s.stopped:
			return nil
		case <-time.After(15 * time.Minute):

		}

	}
}

func (s *scheduler) stop(err error) {
	s.stopped <- struct{}{}
}
