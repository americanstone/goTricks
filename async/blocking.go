package main
func main(){
	chanA :=make(chan string)
	chanB :=make(chan string)
	go func(){
		chanA <- "e1"
		chanB <- "e2"

	}()
	// go func(){

	// 	chanB <- "e2"
	// }()
	<- chanB
	<- chanA
}
