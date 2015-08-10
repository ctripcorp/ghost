package retry

import (
	"time"
)

//Recursion
var (
	Linear = func(fir, pre time.Duration) (cur time.Duration) { return pre + fir }
	Double = func(fir, pre time.Duration) (cur time.Duration) { return pre * 2 }
)

//Transform
var (
	Max = func(max time.Duration) (ret Transform) {
		return func(ori time.Duration) (ret time.Duration) {
			if ori > max {
				return max
			} else {
				return ori
			}
		}
	}
)
