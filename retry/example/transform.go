package main

import (
	retry "../"
	"errors"
	"fmt"
	"time"
)

func main() {
	r := retry.Retry{

		FirstSleep: 1 * time.Second,

		Recursion: retry.Double,

		//the func to limit or randomize sleep time
		//example:
		//
		//Transform: 
		//	func(ori time.Duration)(ret time.Duration) {
		//		if ori > 5 * time.Second {
		//			return 5 * time.Second
		//		} else {
		//			return ori
		//		}
		//	},
		Transform: retry.Max(3 * time.Second),

		Retries: 5,
	}
	r.Attempt(func() error {
		fmt.Println(time.Now().Format(time.Stamp))
		return errors.New("error")
	})
}
