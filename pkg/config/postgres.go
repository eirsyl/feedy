package config

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/spf13/viper"

	"github.com/lib/pq"
)

type postgresConfig struct {
	db *sql.DB
}

// NewPostgresConfig creates a new config backend where postgres is used for storage
func NewPostgresConfig() (Config, error) {

	var host, user, password, database string
	var port int64

	{
		host = viper.GetString("postgresHost")
		user = viper.GetString("postgresUser")
		password = viper.GetString("postgresPassword")
		database = viper.GetString("postgresDatabase")
		port = viper.GetInt64("postgresPort")
	}

	connection := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		database,
	)

	db, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	c := &postgresConfig{
		db: db,
	}

	if err = c.migrate(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *postgresConfig) SaveUser(user User) error {

	var statement = `
	INSERT INTO users
		(
			id,
			consumerKey,
			token
		)
	VALUES
		(
			$1,
			$2,
			$3
		)
	ON CONFLICT (id)
	DO
		UPDATE
			SET consumerKey = $2, token = $3;
	`

	rows, err := c.db.Query(statement, 1, user.ConsumerKey, user.Token)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (c *postgresConfig) GetUser() (User, error) {

	var user User

	var statement = `
	SELECT
		consumerKey, token
	FROM
		users
	WHERE
		id = $1;
	`

	err := c.db.QueryRow(statement, 1).Scan(&user.ConsumerKey, &user.Token)
	return user, err
}

func (c *postgresConfig) GetFeeds() ([]Feed, error) {
	var feeds []Feed

	var statement = `
	SELECT
		name, url, tags
	FROM
		feeds;
	`

	rows, err := c.db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var feed Feed

		if err = rows.Scan(&feed.Name, &feed.URL, pq.Array(&feed.Tags)); err != nil {
			return nil, err
		}

		feeds = append(feeds, feed)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return feeds, nil
}

func (c *postgresConfig) AddFeed(feed Feed) error {

	var statement = `
	INSERT INTO feeds
		(
			name,
			url,
			tags
		)
	VALUES
		(
			$1,
			$2,
			$3
		)
	ON CONFLICT (url)
	DO
		UPDATE
			SET name = $1, tags = $3;
	`

	rows, err := c.db.Query(statement, feed.Name, feed.URL, pq.Array(feed.Tags))
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (c *postgresConfig) RemoveFeed(url string) error {

	var statement = `
	DELETE FROM feeds
	WHERE
		url = $1;
	`

	result, err := c.db.Exec(statement, url)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrFeedNotFound
	}

	return nil
}

func (c *postgresConfig) AddScrapedURL(url string) error {

	var statement = `
	INSERT INTO scrapes
		(
			url
		)
	VALUES
		(
			$1
		)
	ON CONFLICT DO NOTHING;
	`

	rows, err := c.db.Query(statement, strings.ToLower(url))
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (c *postgresConfig) IsScrapedURL(url string) (bool, error) {

	var exists bool

	var statement = `
	SELECT EXISTS (
		SELECT
			url
		FROM
			scrapes
		WHERE
			url = $1
	);
	`

	err := c.db.QueryRow(statement, strings.ToLower(url)).Scan(&exists)

	return exists, err
}

func (c *postgresConfig) Close() error {
	return c.db.Close()
}

/*
* Helpers
 */

// migrate creates the database tables requred for information storage
func (c *postgresConfig) migrate() error {

	var statement = `
	CREATE TABLE IF NOT EXISTS users (
		id serial PRIMARY KEY,
		consumerKey text NOT NULL,
		token text NOT NULL
	);

	CREATE TABLE IF NOT EXISTS feeds (
		name text NOT NULL,
		url text NOT NULL PRIMARY KEY,
		tags text[]
	);

	CREATE TABLE IF NOT EXISTS scrapes (
		url text NOT NULL PRIMARY KEY
	);
	`

	_, err := c.db.Exec(statement)
	return err
}
