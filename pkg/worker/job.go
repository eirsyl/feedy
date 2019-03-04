package worker

import (
	"time"

	"github.com/eirsyl/feedy/pkg/config"
)

type job struct {
	feed config.Feed

	createTime *time.Time
	startTime  *time.Time
}

type jobResult struct {
	job *job
	err error
}
