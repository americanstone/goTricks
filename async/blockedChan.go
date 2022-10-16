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

	c := make(chan urlStatus)
	for _, url := range urls {
		go checkUrl(url, c)
	}
	//Receivers always block until there is data to receive
	//https://go.dev/doc/effective_go#channels
	for r := range c {
		if r.status {
			fmt.Println(r.url, "is up.")
		} else {
			fmt.Println(r.url, "is down!!")
		}
	}   
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
