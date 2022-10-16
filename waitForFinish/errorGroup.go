package main

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	var group errgroup.Group

	greet := func() error {
		time.Sleep(time.Second)
		fmt.Println("Hello 👋")
		return nil
	}

	group.Go(greet)

	// get the first error from the group
	if err := group.Wait(); err != nil {
		log.Fatal("We have an error:", err)
	}

}
