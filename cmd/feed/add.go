package feed

import (
	"fmt"
	"net/url"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/eirsyl/feedy/pkg/scraper"
)

var addCmd = &cobra.Command{
	Use:   "add [feed]",
	Short: "Add feed",
	Long: `
Add feed to the list of feeds to scrape.
	`,
	Args: cobra.ExactArgs(1),
	PreRun: func(_ *cobra.Command, args []string) {
	},
	RunE: func(_ *cobra.Command, args []string) error {

		feedURL := args[0]

		u, err := url.Parse(feedURL)
		if err != nil {
			return errors.Wrap(err, "could not parse feed url")
		}

		s, err := scraper.New()
		if err != nil {
			return errors.Wrap(err, "could not initialize scraper")
		}

		meta, err := s.DiscoverFeed(u)
		if err != nil {
			return errors.Wrap(err, "could not load feed metadata")
		}

		fmt.Println("Title: ", meta.Title)
		fmt.Println("Description: ", meta.Description)
		fmt.Println("Author: ", meta.Author)

		return nil
	},
}
