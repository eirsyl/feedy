package worker

import (
	"net/url"
	"time"

	"github.com/eirsyl/feedy/pkg/pocket"

	"github.com/eirsyl/feedy/pkg/config"

	"github.com/eirsyl/feedy/pkg/scraper"

	"github.com/pkg/errors"
)

type processor struct {
	jobs    <-chan *job
	results chan<- *jobResult

	user    config.User
	config  config.Config
	scraper scraper.Scraper
	pocket  pocket.Pocket
}

func newProcessor(
	jobs <-chan *job,
	results chan<- *jobResult,

	user config.User,
	config config.Config,
	scraper scraper.Scraper,
	pocket pocket.Pocket,
) (*processor, error) {
	p := &processor{
		jobs:    jobs,
		results: results,

		user:    user,
		config:  config,
		scraper: scraper,
		pocket:  pocket,
	}

	return p, nil
}

func (p *processor) run() {
	for {

		j, more := <-p.jobs
		if more {
			err := p.process(j)
			p.results <- &jobResult{job: j, err: err}
		} else {
			return
		}

	}
}

func (p *processor) process(j *job) error {
	//ctx := context.TODO()

	var currentTime = time.Now()
	j.startTime = &currentTime

	// Parse url
	u, err := url.Parse(j.feed.URL)
	if err != nil {
		return errors.Wrapf(err, "invalid feed url %s", j.feed.URL)
	}

	results, err := p.scraper.ScrapeFeed(u)
	if err != nil {
		return errors.Wrapf(err, "could not scrape feed %s", j.feed.URL)
	}

	for _, result := range results {

		isScraped, err := p.config.IsScrapedURL(result.URL)
		if err != nil {
			return errors.Wrapf(err, "could not lookup scrape status for url %s", result.URL)
		}

		if !isScraped {
			if err = p.pocket.AddItem(result.URL, result.Name, j.feed.Tags, p.user.ConsumerKey, p.user.Token); err != nil {
				return errors.Wrapf(err, "could not add %s to pocket", result.URL)
			}

			if err = p.config.AddScrapedURL(result.URL); err != nil {
				return errors.Wrapf(err, "could not mark url %s as scraped", result.URL)
			}
		}
	}

	return nil
}
