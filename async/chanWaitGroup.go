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

/*
 one of way to tell async jobs are completed via waitGroup
 don't like this pattern since the execution order in main and
 goroutine is random
*/
func main() {
	c := make(chan urlStatus)

	wg := &sync.WaitGroup{}

	for _, url := range urls {
		wg.Add(1)
		go checkUrl(url, c, wg)
	}
	fmt.Println("Main: async check Urls")

	// this pattern helps the performance since you can do
	// other async jobs before main take the flow control
	go consumer(c, wg)

	fmt.Println("Main: async print")

	fmt.Println("Main: Waiting for complemete!")
	wg.Wait() //only wait for checkUrls complete
	close(c)
	fmt.Println("Main: Wait complemete!") // this print first then the consumer prints
}
func consumer(c <-chan urlStatus, wg *sync.WaitGroup) {
	// this will block until close(c) is called in main control
	for r := range c {
		if r.status {
			fmt.Println(r.url, "is up.")
		} else {
			fmt.Println(r.url, "is down!!")
		}
	}
	fmt.Println("end of consumer")
}

//checks and prints a message if a website is up or down
func checkUrl(url string, c chan<- urlStatus, wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := http.Get(url)
	if err != nil {
		// The website is down
		c <- urlStatus{url, false}
	} else {
		// The website is up
		c <- urlStatus{url, true}
	}
}

type urlStatus struct {
	url    string
	status bool
}
