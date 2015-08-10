package main

import (
	retry "../"
	"errors"
	"fmt"
	"time"
)

func main() {
	//retry at most 3 times.
	//sleeps for 1 second before each retry
	retries, errors := retry.Attempt(3, 1*time.Second, func() error {
		err := errors.New("myError")
		println(err.Error())
		return err
	})
	//retries = 3
	//errors = [myError myError myError myError]
	fmt.Printf("retries = %d\nerrors = %v\n", retries, errors)
}
