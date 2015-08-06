package retry

import (
	"testing"
	"time"
	"errors"
	"fmt"
)

func TestDefaultRetry (t *testing.T) {
	retries, errors := Attemp(3, 1*time.Second, func()error{
		err := errors.New("Error counted, please retry.")
		println(err.Error())
		return err
	})
	println("Retry time: ", retries)
	fmt.Printf("%v\n", errors)
}
