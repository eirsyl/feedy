package feed

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jedib0t/go-pretty/text"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/eirsyl/feedy/pkg/client"
	"github.com/eirsyl/feedy/pkg/config"
	"github.com/eirsyl/feedy/pkg/scraper"
)

var addCmd = &cobra.Command{
	Use:   "add [feed] [...tags]",
	Short: "Add feed",
	Long: `
Add feed to the list of feeds to scrape.
	`,
	Args: cobra.MinimumNArgs(1),
	PreRun: func(_ *cobra.Command, args []string) {
	},
	RunE: func(_ *cobra.Command, args []string) error {

		feedURL := args[0]
		tags := args[1:]

		c, err := config.GetConfigProvider()
		if err != nil {
			return err
		}
		defer c.Close() // nolint: errcheck

		u, err := url.Parse(feedURL)
		if err != nil {
			return errors.Wrap(err, "could not parse feed url")
		}

		cc, err := client.New()
		if err != nil {
			return errors.Wrap(err, "could not create client")
		}

		s, err := scraper.New(cc)
		if err != nil {
			return errors.Wrap(err, "could not initialize scraper")
		}

		meta, err := s.DiscoverFeed(context.TODO(), u)
		if err != nil {
			return errors.Wrap(err, "could not load feed metadata")
		}

		if err = c.AddFeed(config.Feed{
			URL:  meta.URL,
			Name: meta.Title,
			Tags: tags,
		}); err != nil {
			return errors.Wrap(err, "could not add feed to config")
		}

		// Mark existing elements as scraped de dont want to fill up pocket when we add a new feed
		for _, item := range meta.Items {
			err = c.AddScrapedURL(item)
			if err != nil {
				return errors.Wrap(err, "could not add item to the list of scraped urls")
			}
		}

		fmt.Println(text.FgGreen.Sprintf("Added %s (%s) to the watch list", meta.Title, meta.URL))

		return nil
	},
}
