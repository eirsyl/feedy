package config

type fileConfig struct {
}

// NewFileConfig returns a new config backend based on a file (yaml)
func NewFileConfig() (Config, error) {
	return &fileConfig{}, nil
}

func (c *fileConfig) SaveToken(token string, overwrite bool) error {
	return nil
}

func (c *fileConfig) GetToken() (string, error) {
	return "", nil
}

func (c *fileConfig) GetFeeds() ([]Feed, error) {
	return []Feed{}, nil
}

func (c *fileConfig) SaveFeed(feed Feed) error {
	return nil
}
