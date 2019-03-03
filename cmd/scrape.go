package cmd

import (
	"net/url"

	"github.com/eirsyl/feedy/pkg/client"
	"github.com/eirsyl/feedy/pkg/config"
	"github.com/eirsyl/feedy/pkg/pocket"
	"github.com/eirsyl/feedy/pkg/scraper"
	"github.com/eirsyl/flexit/log"
	"github.com/pkg/errors"
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

		logger := log.NewLogrusLogger(false)

		c, err := config.NewFileConfig()
		if err != nil {
			return errors.Wrap(err, "could not create config backend")
		}
		defer c.Close()

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

		for _, feed := range feeds {
			logger.Infof("Scraping feed %s:", feed.Name)

			u, err := url.Parse(feed.URL)
			if err != nil {
				return errors.Wrapf(err, "invalid feed url %s", feed.URL)
			}

			results, err := scr.ScrapeFeed(u)
			if err != nil {
				return errors.Wrapf(err, "could not scrape feed %s", feed.Name)
			}

			for _, result := range results {

				isScraped, err := c.IsScrapedURL(result.URL)
				if err != nil {
					return errors.Wrapf(err, "could not lookup scrape status for url %s", result.URL)
				}

				if !isScraped {
					if err = pocket.AddItem(result.URL, result.Name, feed.Tags, user.ConsumerKey, user.Token); err != nil {
						return errors.Wrapf(err, "could not add %s to pocket", result.URL)
					}

					if err = c.AddScrapedURL(result.URL); err != nil {
						return errors.Wrapf(err, "could not mark url %s as scraped", result.URL)
					}
				}
			}
		}

		return nil
	},
}
