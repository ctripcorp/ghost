package main

import (
	retry "../"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	r := retry.Retry{
		//time to sleep before first retry
		FirstSleep: 1 * time.Second,

		//the func to compute next sleep time
		//example:
		//
		//Recursion:
		//	func(fir, pre time.Duration) (cur time.Duration) {
		//		return pre + fir
		//	},
		Recursion: retry.Linear,

		//max time to retry
		Retries: 3,
	}
	//return nil to announce success
	//return a counted error to announce failure requesting for retry
	r.Attempt(func() error {
		fmt.Println(time.Now().Format("+ " + time.Stamp))
		resp, err := http.Get("http://www.facebook.com/")
		fmt.Println(time.Now().Format("- " + time.Stamp))
		if err != nil {
			//return err to invoke retry
			return err
		}
		defer resp.Body.Close()
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			//return err to invoke retry
			return err
		}
		//return nil to announce success
		return nil
	})
}
