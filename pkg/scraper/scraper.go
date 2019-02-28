package scraper

import (
	"net/url"

	"github.com/mmcdole/gofeed"
)

// FeedMeta contains metadata about a feed
type FeedMeta struct {
	Title       string
	Description string
	Author      string
}

// Result contains a single feed element
type Result struct {
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
func New() (Scraper, error) {

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

	return &FeedMeta{
		Title:       feed.Title,
		Description: feed.Description,
		Author:      author,
	}, nil
}

func (s *baseScraper) ScrapeFeed(u *url.URL) ([]Result, error) {
	return []Result{}, nil
}
