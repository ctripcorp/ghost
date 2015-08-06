package retry

import (
	"time"
)

var (
	Double = func(bSub1 time.Duration) (n time.Duration) { return bSub1 * 2 }
)

type Operation func() error

type Recursion func(nSub1 time.Duration) (n time.Duration)

type Retry interface {
	Attemp(op Operation) (int, []error)
}

//a struct to store retry strategy params
type retry struct {
	//switch of randomizing
	randomize bool

	//sleep time before first retry
	firstSleep time.Duration
	//range of sleep time
	minSleep time.Duration
	maxSleep time.Duration

	//a func to compute sleep time
	recursion Recursion

	//max retry time
	retries int
}

func (r retry) Attemp(op Operation) (retries int, errors []error) {
	retries = 0
	errors = nil
	err := op()
	if err == nil {
		return retries, errors
	}
	errors = make([]error, 0)
	sleep := r.firstSleep
	for retries < r.retries {
		time.Sleep(r.limit(sleep))
		retries++
		err = op()
		if err != nil {
			sleep = r.recursion(sleep)
			errors = append(errors, err)
		} else {
			break
		}
	}
	return retries, errors
}

func (r retry) limit(d time.Duration) time.Duration {
	if d > r.maxSleep {
		return r.maxSleep
	}
	if d < r.minSleep {
		return r.minSleep
	}
	return d
}

func Attemp(retries int, firstSleep time.Duration, op Operation) (int, []error){

	r := &retry{
		randomize:  false,
		firstSleep: firstSleep,
		minSleep:   0 * time.Second,
		maxSleep:   60 * time.Second,
		recursion:  Double,
		retries:    retries,
	}

	return r.Attemp(op)
}
