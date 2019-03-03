package config

// Config defines the required methods for a config backend
type Config interface {
	// Authentication management
	SaveUser(user User) error
	GetUser() (User, error)
	// Feed management
	GetFeeds() ([]Feed, error)
	AddFeed(feed Feed) error
	RemoveFeed(url string) error
	// Scrape management
	AddScrapedURL(url string) error
	IsScrapedURL(url string) (bool, error)
	// Close config provider after use
	Close() error
}

// Feed defines the feed structure stored on disk
type Feed struct {
	URL  string   `json:"url"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

// User stores the user credential used when contacting pocket
type User struct {
	ConsumerKey string `json:"consumerKey"`
	Token       string `json:"token"`
}

// Validate validates the feed content
func (f *Feed) Validate() error {
	return nil
}
