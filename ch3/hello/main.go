package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go SayHello(&wg)
	wg.Wait()
}

func SayHello(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Hello")
}
