package retry

import (
	"time"
)

type Operation func() error

type Transform func(ori time.Duration) (ret time.Duration)

type Recursion func(fir, pre time.Duration) (cur time.Duration)

//a struct to store retry strategy params
type Retry struct {
	//sleep time before first retry
	FirstSleep time.Duration

	//a func to adjust(limit, randomize or tamper) current sleep time
	Transform Transform

	//a func to compute next sleep time
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
	fir := r.FirstSleep
	sleep := r.FirstSleep
	for retries < r.Retries {
		if r.Transform != nil {
			time.Sleep(r.Transform(sleep))
		} else {
			time.Sleep(sleep)
		}
		retries++
		err = op()
		if err != nil {
			if r.Recursion != nil {
				sleep = r.Recursion(fir, sleep)
			}
			errors = append(errors, err)
		} else {
			break
		}
	}
	return retries, errors
}
