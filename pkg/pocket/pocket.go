package pocket

import (
	"context"
	"strings"

	"github.com/eirsyl/feedy/pkg/client"
	"github.com/pkg/errors"
)

// Pocket defines the methods used when onteracting with pocket
type Pocket interface {
	Login(consumerKey string) (string, error)
	AddItem(url, name string, tags []string, consumerKey, token string) error
}

type basePocket struct {
	c *client.Client
}

// New creates a new instance of the pocket client
func New(c *client.Client) (Pocket, error) {
	return &basePocket{
		c: c,
	}, nil
}

func (p *basePocket) AddItem(url, name string, tags []string, consumerKey, token string) error {
	addItemBody, err := p.c.NewRequest("POST", "https://getpocket.com/v3/add", addItemRequest{
		URL:         url,
		Title:       name,
		Tags:        strings.Join(tags, ","),
		ConsumerKey: consumerKey,
		AccessToken: token,
	})
	if err != nil {
		return errors.Wrap(err, "could not create add item request")
	}

	_, err = p.c.Do(context.TODO(), addItemBody, nil)
	if err != nil {
		return errors.Wrap(err, "could not add item")
	}

	return nil
}

type addItemRequest struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Tags        string `json:"tags"`
	ConsumerKey string `json:"consumer_key"`
	AccessToken string `json:"access_token"`
}
