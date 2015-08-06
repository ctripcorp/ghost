package main

import (
	retry "../"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	r := &retry.Retry{
		//switch to activate randomize
		Randomize: false,

		//time to sleep before first retry
		FirstSleep: 1 * time.Second,
		//range of sleep time
		MinSleep: 0 * time.Second,
		MaxSleep: 3 * time.Second,

		//sleep time increase or decline
		//linear:
		//	n = nsub1 + 1
		//exponent:
		//	n = nsub1 * 2
		Recursion: retry.Double,
		//Recursion: func(nsub1 time.Duration) time.Duration {
		//	return nsub1+1 * time.Second
		//},

		//max time to retry
		Retries: 5,
	}
	//return nil to announce success
	//return a counted error to announce failure, requesting for retry
	r.Attempt(func() error {
		resp, err := http.Get("http://www.baidu.com/")
		//resp, err := http.Get("http://www.facebook.com/")
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("return err to invoke retry")
			return err
		}
		defer resp.Body.Close()
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("return err to invoke retry")
			return err
		}
		fmt.Println("return nil to announce success")
		return nil
	})
}
