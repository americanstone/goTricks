package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var lock sync.Mutex
	greet := func() {
		defer lock.Unlock()
		time.Sleep(time.Second)
		fmt.Println("hey there")
	}
	lock.Lock()

	go greet()

	lock.Lock()
}
