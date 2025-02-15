package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	newRandStream := func(
		done <-chan struct{},
	) <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited")
			defer close(randStream)
			for {
				select {
				case randStream <- rand.Int():
				case <-done:
					return
				}
			}
		}()
		return randStream
	}

	done := make(chan struct{})
	randStream := newRandStream(done)
	fmt.Println("3 random ints")
	for i := 0; i < 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
	close(done)
	time.Sleep(1 * time.Second)
}
