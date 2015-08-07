package retry

import (
	"time"
)

func Attempt(retries int, firstSleep time.Duration, op Operation) (int, []error){

	r := &Retry{
		Randomize:  false,
		FirstSleep: firstSleep,
		MinSleep:   0 * time.Second,
		MaxSleep:   60 * time.Second,
		Recursion:  Double,
		Retries:    retries,
	}

	return r.Attempt(op)
}
