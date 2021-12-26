package main

import (
	"fmt"
	"main/prime"
)

func main() {
	p := prime.Primes(1000000)
	fmt.Println(p)
}
