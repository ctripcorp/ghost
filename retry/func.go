package retry

import (
	"time"
)

func Attempt(retries int, firstSleep time.Duration, op Operation) (int, []error) {

	r := &Retry{
		FirstSleep: firstSleep,
		Retries:    retries,
	}

	return r.Attempt(op)
}
