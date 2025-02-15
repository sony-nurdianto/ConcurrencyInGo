package main

import (
	"fmt"
	"time"
)

func main() {
	dowork := func(
		done <-chan struct{},
		s <-chan string,
	) <-chan interface{} {
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited")
			defer close(terminated)
			for {
				select {
				case w := <-s:
					// dosomething
					fmt.Println(w)
				case <-done:
					return
				}
			}
		}()

		return terminated
	}

	done := make(chan struct{})
	terminated := dowork(done, nil)

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling dowork goroutine")
		close(done)
	}()

	<-terminated
	fmt.Println(done)
}
