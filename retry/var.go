package retry

import (
	"time"
)

var (
	Double = func(nSub1 time.Duration) (n time.Duration) { return nSub1 * 2 }
)
