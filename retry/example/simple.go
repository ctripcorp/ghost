package main

import (
	"errors"
	"fmt"
	"time"
	retry "../"
)

func main() {
	//Retry at most 3 times.
	//Sleeps for 1 second before first retry, and sleep time doubles after each time it retries
	retries, errors := retry.Attempt(3, 1*time.Second, func() error {
		err := errors.New("myError")
		println(err.Error())
		return err
	})
	//retries = 3
	//errors = [myError myError myError myError]
	fmt.Printf("retries = %d\nerrors = %v\n", retries, errors)
}
