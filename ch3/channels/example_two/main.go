package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	ch := make(chan byte, 3)

	words := []byte("ABCD")

	wg.Add(len(words))

	for _, v := range words {
		go func(b byte) {
			defer wg.Done()
			ch <- b
		}(v)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for b := range ch {
		fmt.Println(string(b))
	}
}
