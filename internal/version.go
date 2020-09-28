package internal

import "fmt"

// Version stores the application version
var Version string

// BuildDate contains the build date
var BuildDate string

// ShortDescription contains a short description of this application
var ShortDescription = "Feedy - RSS feed scraper"

// LongDescription contains a longer description of this application
var LongDescription = fmt.Sprintf(
	`%s

Scrape RSS feeds and send the articles to getpocket.com.
`,
	ShortDescription,
)
