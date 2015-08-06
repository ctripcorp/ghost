package retry

import (
	"time"
)

var (
	Double = func(nSub1 time.Duration) (n time.Duration) { return nSub1 * 2 }
)

type Operation func() error

type Recursion func(nSub1 time.Duration) (n time.Duration)

//a struct to store retry strategy params
type Retry struct {
	//switch of randomizing
	Randomize bool

	//sleep time before first retry
	FirstSleep time.Duration
	//range of sleep time
	MinSleep time.Duration
	MaxSleep time.Duration

	//a func to compute sleep time
	Recursion Recursion

	//max retry time
	Retries int
}

func (r Retry) Attempt(op Operation) (retries int, errors []error) {
	retries = 0
	errors = nil
	err := op()
	if err == nil {
		return retries, errors
	}
	errors = make([]error, 0)
	errors = append(errors, err)
	sleep := r.FirstSleep
	for retries < r.Retries {
		time.Sleep(r.limit(sleep))
		retries++
		err = op()
		if err != nil {
			sleep = r.Recursion(sleep)
			errors = append(errors, err)
		} else {
			break
		}
	}
	return retries, errors
}

func (r Retry) limit(d time.Duration) time.Duration {
	if d > r.MaxSleep {
		return r.MaxSleep
	}
	if d < r.MinSleep {
		return r.MinSleep
	}
	return d
}

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
