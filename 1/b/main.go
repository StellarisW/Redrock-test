package main

import "fmt"

func main() {
	ch := make(chan int)
	var receive int
	go func() {
		ch <- 1
		fmt.Println("下山的路又堵起了")
	}()
	receive = <-ch
	fmt.Println(receive)
}
