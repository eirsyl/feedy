package runtime

import (
	"github.com/eirsyl/feedy/pkg/utils/log"
	"go.uber.org/automaxprocs/maxprocs"
)

// OptimizeRuntime configures the runtime based on available resources.
func OptimizeRuntime(logger log.Logger) error {
	_, err := maxprocs.Set(maxprocs.Logger(logger.Debugf))
	return err
}
