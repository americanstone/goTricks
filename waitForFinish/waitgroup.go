package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	greet := func() {
		defer wg.Done() // we unlock when done

		time.Sleep(time.Second)
		fmt.Println("Hello 👋")
	}

	// how many functions will we call?
	wg.Add(1)

	go greet()

	wg.Wait() //
}
