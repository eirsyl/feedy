package config

// Config defines the required methods for a config backend
type Config interface {
	// Authentication management
	SaveToken(token string, overwrite bool) error
	GetToken() (string, error)
	// Feed management
	GetFeeds() ([]Feed, error)
	SaveFeed(feed Feed) error
}

// Feed defines the feed structure stored on disk
type Feed struct {
	URL  string
	Name string
	Tags []string
}

// Validate validates the feed content
func (f *Feed) Validate() error {
	return nil
}
