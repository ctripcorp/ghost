package retry

import (
	"time"
)

var (
	Double = func(fir, pre time.Duration) (cur time.Duration) { return pre * 2 }
)
