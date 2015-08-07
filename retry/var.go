package retry

import (
	"time"
)

var (
	Linear = func(fir, pre time.Duration) (cur time.Duration) { return pre + fir }
	Double = func(fir, pre time.Duration) (cur time.Duration) { return pre * 2 }
)
