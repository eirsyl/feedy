package worker

import (
	"sync"
	"time"

	"github.com/eirsyl/flexit/log"

	"github.com/eirsyl/feedy/pkg/config"
	"github.com/eirsyl/feedy/pkg/pocket"
	"github.com/eirsyl/feedy/pkg/scraper"
)

// Worker defines the worker interface used to scrape multiple feeds at once
type Worker interface {
	Add(feed config.Feed) error
	Run() error
	Stop(error)
}

// nolint: maligned
type baseWorker struct {
	concurrency int
	autostop    bool

	jobs    chan *job
	results chan *jobResult

	stop    chan struct{}
	stopped bool
	wg      sync.WaitGroup

	user    config.User
	config  config.Config
	scraper scraper.Scraper
	pocket  pocket.Pocket
	logger  log.Logger
}

// New returns a new instance of the worker
func New(
	concurrency int,
	autostop bool,

	user config.User,
	config config.Config,
	scraper scraper.Scraper,
	pocket pocket.Pocket,
	logger log.Logger,
) (Worker, error) {
	worker := &baseWorker{
		concurrency: concurrency,
		autostop:    autostop,

		jobs:    make(chan *job, 100),
		results: make(chan *jobResult, 100),

		stop:    make(chan struct{}, 1),
		stopped: false,
		wg:      sync.WaitGroup{},

		user:    user,
		config:  config,
		scraper: scraper,
		pocket:  pocket,
		logger:  logger,
	}

	return worker, nil
}

func (w *baseWorker) Add(feed config.Feed) error {
	if w.stopped {
		return ErrWorkerClosed
	}

	var createTime = time.Now()

	job := &job{
		feed:       feed,
		createTime: &createTime,
	}

	w.jobs <- job
	w.wg.Add(1)

	return nil
}

func (w *baseWorker) Run() error {
	var err error

	// Start background workers
	for id := 1; id <= w.concurrency; id++ {
		var p *processor

		p, err = newProcessor(
			w.jobs,
			w.results,
			w.user,
			w.config,
			w.scraper,
			w.pocket,
		)
		if err != nil {
			return err
		}

		go p.run()
	}

	// Start result collection
	go func() {
		for {
			r, more := <-w.results
			currentTime := time.Now()

			if more {
				if r.err != nil {
					w.logger.WithField("feed", r.job.feed.URL).Errorf(
						"Job failed with error: %v", r.err,
					)
				} else {
					w.logger.WithField("feed", r.job.feed.URL).Infof(
						"Job succeeded, runtime: %s", currentTime.Sub(*r.job.createTime),
					)
				}

				w.wg.Done()
			} else {
				return
			}
		}
	}()

	// Wait on close signal and then all pending builds
	if !w.autostop {
		<-w.stop
	}
	w.wg.Wait()

	// Close jobs and results channels, this stops the workers
	close(w.jobs)
	close(w.results)

	return nil
}

func (w *baseWorker) Stop(err error) {
	w.stop <- struct{}{}
}
