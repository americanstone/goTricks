// https://scene-si.org/2020/05/29/waiting-on-goroutines/
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	greet := func() {
		defer cancel()
		time.Sleep(time.Second)
		fmt.Println("hey there")
	}
	go greet()
	<-ctx.Done()
}
