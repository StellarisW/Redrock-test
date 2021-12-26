package main

import "fmt"

func printer(ch chan int) {
	for {
		data := <-ch
		if data == 0 {
			break
		}
		fmt.Println("下山的路又堵起了")
	}
	ch <- 0
}

func main() {
	ch := make(chan int)
	go printer(ch)
	for i := 1; i <= 10; i++ {
		ch <- i
	}
	ch <- 0
	<-ch
}
