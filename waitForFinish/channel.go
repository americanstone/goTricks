package main

import (
	"fmt"
	"time"
)

func main() {
	finished := make(chan bool)
	greet := func() {
		time.Sleep(time.Second)
		fmt.Println("Hey there")
		finished <- true
	}
	go greet()
	<-finished
}
