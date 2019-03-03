package config

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	bolt "go.etcd.io/bbolt"
)

var (
	authenticationBucket = []byte("Authentication")
	feedBucket           = []byte("Feed")
	scrapeBucket         = []byte("Scrapes")
)

type fileConfig struct {
	db *bolt.DB
}

// NewFileConfig returns a new config backend based on a file (yaml)
func NewFileConfig() (Config, error) {
	path := viper.GetString("configFile")
	if path == "" {
		return nil, ErrConfigFileNotGiven
	}

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not open database")
	}

	c := &fileConfig{
		db: db,
	}

	if err = c.migrate(); err != nil {
		return nil, errors.Wrap(err, "could not migrate database")
	}

	return c, nil
}

func (c *fileConfig) SaveUser(user User) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(authenticationBucket)

		buf, err := json.Marshal(user)
		if err != nil {
			return errors.Wrap(err, "could not marshal user")
		}

		return b.Put([]byte("default"), buf)
	})
}

func (c *fileConfig) GetUser() (User, error) {
	var u User

	if err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(authenticationBucket)

		buf := b.Get([]byte("default"))

		if err := json.Unmarshal(buf, &u); err != nil {
			return errors.Wrap(err, "could not unmarshal user")
		}

		return nil
	}); err != nil {
		return User{}, errors.Wrap(err, "could not read from database")
	}

	return u, nil
}

func (c *fileConfig) GetFeeds() ([]Feed, error) {
	var feeds []Feed

	if err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(feedBucket)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			var f Feed

			if err := json.Unmarshal(v, &f); err != nil {
				return err
			}

			feeds = append(feeds, f)
		}

		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "could not iterate over feeds")
	}

	return feeds, nil
}

func (c *fileConfig) AddFeed(feed Feed) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(feedBucket)

		buf, err := json.Marshal(feed)
		if err != nil {
			return errors.Wrap(err, "could not marshal feed")
		}

		return b.Put([]byte(feed.URL), buf)
	})
}

func (c *fileConfig) RemoveFeed(url string) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(feedBucket)

		return b.Delete([]byte(url))
	})
}

func (c *fileConfig) AddScrapedURL(url string) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(scrapeBucket)

		return b.Put([]byte(url), []byte("1"))
	})
}

func (c *fileConfig) IsScrapedURL(url string) (bool, error) {

	var scraped bool

	if err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(scrapeBucket)

		scraped = string(b.Get([]byte(url))) == "1"

		return nil
	}); err != nil {
		return false, errors.Wrap(err, "could not lookup scrape status")
	}

	return scraped, nil
}

func (c *fileConfig) Close() error {
	return c.db.Close()
}

func (c *fileConfig) migrate() error {
	return c.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(authenticationBucket)
		if err != nil {
			return errors.Wrap(err, "could not create bucket")
		}

		_, err = tx.CreateBucketIfNotExists(feedBucket)
		if err != nil {
			return errors.Wrap(err, "could not create bucket")
		}

		_, err = tx.CreateBucketIfNotExists(scrapeBucket)
		if err != nil {
			return errors.Wrap(err, "could not create bucket")
		}

		return nil
	})
}
