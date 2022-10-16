package main

import (
	"fmt"
	"time"
)

var semaphore = make(chan int, 5)

func main(){
	for i:=0; i<= 1000; i++{
		go handle(i)
	}
	fmt.Printf("Done \n")
	fmt.Scanln()
}

func handle(i int){
	semaphore <- 1
	process(i)
	<-semaphore
}
func process(i int){
	fmt.Printf("Processing %v\n", i)
	time.Sleep(time.Second)
	fmt.Printf("Processed %v\n",i)
}