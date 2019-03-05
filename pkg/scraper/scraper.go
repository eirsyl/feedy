package scraper

import (
	"context"
	"net/url"
	"strings"

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
	DiscoverFeed(ctx context.Context, u *url.URL) (*FeedMeta, error)
	ScrapeFeed(ctx context.Context, u *url.URL) ([]Result, error)
}

type baseScraper struct {
	parser *gofeed.Parser
	client *client.Client
}

// New returns a new scraper instance
func New(c *client.Client) (Scraper, error) {

	parser := gofeed.NewParser()

	return &baseScraper{
		parser: parser,
		client: c,
	}, nil
}

func (s *baseScraper) DiscoverFeed(ctx context.Context, u *url.URL) (*FeedMeta, error) {

	feed, err := s.scrape(ctx, u)
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

func (s *baseScraper) ScrapeFeed(ctx context.Context, u *url.URL) ([]Result, error) {

	var result []Result

	feed, err := s.scrape(ctx, u)
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

func (s *baseScraper) scrape(ctx context.Context, u *url.URL) (*gofeed.Feed, error) {

	req, err := s.client.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/xhtml+xml,application/xml;q=0.9")

	// tumblr.com presents the users with a GDPR screen if the user-agent is present
	if strings.Contains(u.Host, "tumblr.com") {
		req.Header.Del("User-Agent")
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // nolint: gas, errcheck

	return s.parser.Parse(resp.Body)
}
