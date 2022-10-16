package main

import (
	"fmt"
	"net/http"
)

func main() {
	// A slice of sample websites
	urls := []string{
		"https://www.easyjet.com/",
		"https://www.skyscanner.de/",
		"https://www.ryanair.com",
		"https://wizzair.com/",
		"https://www.swiss.com/",
	}
	done := make(chan bool)
	c := producer(urls)

	go consumer(c, done)

	fmt.Println("main: block on complete signal")
	// blocked here because consumer method never put data in done channel
	<-done // block until receieve value
}

func producer(urls []string) chan urlStatus {
	c := make(chan urlStatus)
	for _, url := range urls {
		go checkUrl(url, c)
	}
	return c
}

func consumer(c <-chan urlStatus, done chan<- bool) {
	// print all the data from c and waiting for data arriving until c close
	// sentinel value
	i := 5 // this is a hack I think. this is way to end the loop
	for r := range c {
		if r.status {
			fmt.Println(r.url, "is up.")
		} else {
			fmt.Println(r.url, "is down!!")
		}
		i--
		if 0 == i {
			break
		}
	}
	// so this code will never execute
	fmt.Println("end of consumer loop")
	done <- true
}

//checks and prints a message if a website is up or down
func checkUrl(url string, c chan<- urlStatus) {
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
