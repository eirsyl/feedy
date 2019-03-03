package scraper

import (
	"net/url"

	"github.com/eirsyl/feedy/pkg/client"

	"github.com/mmcdole/gofeed"
)

// FeedMeta contains metadata about a feed
type FeedMeta struct {
	Title       string
	Description string
	Author      string
	URL         string
	Items       []string
}

// Result contains a single feed element
type Result struct {
	URL  string
	Name string
}

// Scraper defines the methods used when interacting with a feed
type Scraper interface {
	DiscoverFeed(u *url.URL) (*FeedMeta, error)
	ScrapeFeed(u *url.URL) ([]Result, error)
}

type baseScraper struct {
	parser *gofeed.Parser
}

// New returns a new scraper instance
func New(c *client.Client) (Scraper, error) {

	parser := gofeed.NewParser()

	return &baseScraper{
		parser: parser,
	}, nil
}

func (s *baseScraper) DiscoverFeed(u *url.URL) (*FeedMeta, error) {

	feed, err := s.parser.ParseURL(u.String())
	if err != nil {
		return nil, err
	}

	var author string
	if feed.Author != nil {
		author = feed.Author.Name
	}

	var items []string
	for _, item := range feed.Items {
		items = append(items, item.Link)
	}

	return &FeedMeta{
		Title:       feed.Title,
		Description: feed.Description,
		Author:      author,
		URL:         u.String(),
		Items:       items,
	}, nil
}

func (s *baseScraper) ScrapeFeed(u *url.URL) ([]Result, error) {

	var result []Result

	feed, err := s.parser.ParseURL(u.String())
	if err != nil {
		return nil, err
	}

	for _, feedItem := range feed.Items {
		result = append(result, Result{
			URL:  feedItem.Link,
			Name: feedItem.Title,
		})
	}

	return result, nil
}
