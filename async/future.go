package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

/*
https://medium.com/strava-engineering/futures-promises-in-the-land-of-golang-1453f4807945#:~:text=The%20future%2Fpromise%20pattern%20is,simple%20implementation%20is%20shown%20below.

	v1 only return data
*/
func RequestFuture(url string) <-chan []byte {
	c := make(chan []byte, 1)
	// goroutine so that return the control to caller
	go func() {
		var body []byte
		resp, err := http.Get(url)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		body, _ = ioutil.ReadAll(resp.Body)
		c <- body
		close(c)
	}()
	return c
}

/*
	v2 return a function blocks on goroutine completion
*/
func RequestFutureFunction(url string) func() ([]byte, error) {
	var body []byte
	var err error
	c := make(chan struct{}, 1)
	go func() {
		defer close(c) // <-- important
		var resp *http.Response
		resp, err = http.Get(url)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
	}()
	// return a function blocks on goroutine completion
	return func() ([]byte, error) {
		<-c // caller will block until defer clost(c) called indicate the async func completed
		return body, err
	}
}

/*

 */
func Future(f func() (interface{}, error)) func() (interface{}, error) {
	var result interface{}
	var err error
	c := make(chan struct{}, 1)
	go func() {
		defer close(c)
		result, err = f()
	}()
	return func() (interface{}, error) {
		<-c
		return result, err
	}
}

// func main() {
// 		future := RequestFuture("http://labs.strava.com")
// 		 // do many other things, maybe create other futures
// 		body := <-future
// 		log.Printf("response length: %d", len(body))
//  }

func main() {
	future := RequestFutureFunction1("http://strava.com")
	// do many other things, maybe create other futures
	body, err := future()
	log.Printf("response length: %d", len(body))
	log.Printf("request error: %v", err)
}

//  func main() {
// 	 url := "http://labs.strava.com"
// 	 future := Future(func() (interface{}, error) {
// 		resp, err := http.Get(url)
// 		if err != nil {
// 			 return nil, err
// 		}
// 		defer resp.Body.Close()
// 		return ioutil.ReadAll(resp.Body)
// 	 })
// 	// do many other things
// 	b, err := future()
// 	body, _ := b.([]byte)
// 	log.Printf("response length: %d", len(body))
// 	log.Printf("request error: %v", err)
//  }
