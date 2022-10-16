package main

import (
	"fmt"
	"net/http"
	"sync"
)

var urls = []string{
	"https://www.easyjet.com/",
	"https://www.skyscanner.de/",
	"https://www.ryanair.com",
	"https://wizzair.com/",
	"https://www.swiss.com/",
}

func checker(urls []string) func() chan urlStatus {
	wg := &sync.WaitGroup{}
	c := make(chan urlStatus)
	for _, url := range urls {
		wg.Add(1)
		go checkUrl(url, c, wg)
	}
	return func() chan urlStatus {
		//because unbuffered channel both sending and retrieving
		//data are blocking. the checkeUrl method sends the data to
		//channel and wait for consumption but my async consumer
		//only starts consuming when waitGrounp is down to 0.
		//since the checker blocked on sending the defer wg.Done()
		// never was executed. so the wg.Wait() will block forever.
		// but if I put wg.Wait() in the goroutine the channel will
		//be returned to my consumer goroutine which unblocks the checkUrl method.
		// so the defer wg.Done() will have the opportunity to run
		// this also can be solved w/o wrapping in goroutine by buffered channel which give the
		// defer wg.Done() to run
		go func() { // why this is necessary
			wg.Wait()
			close(c)
		}()
		return c
	}
}
func main() {
	done := make(chan bool)

	future := checker(urls)

	c := future()

	go consumer(c, done)

	<-done
}

func consumer(c <-chan urlStatus, done chan bool) {
	for r := range c {
		if r.status {
			fmt.Println(r.url, "is up.")
		} else {
			fmt.Println(r.url, "is down!!")
		}
	}
	fmt.Println("end of consumer")
	done <- true
}

func checkUrl(url string, c chan<- urlStatus, wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := http.Get(url)
	if err != nil {
		c <- urlStatus{url, false}
	} else {
		fmt.Println("pushing data to channel")
		c <- urlStatus{url, true}
	}
	fmt.Println("run wg.Done")
}

type urlStatus struct {
	url    string
	status bool
}
