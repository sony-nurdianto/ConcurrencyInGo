package main

import (
	"fmt"
	"math/rand"
	"time"
)

func sendPulse(heartbeat chan<- struct{}) {
	select {
	case heartbeat <- struct{}{}:
	default:
	}
}

func sendResults(
	done <-chan struct{},
	result chan<- time.Time,
	pulse <-chan time.Time,
	heartbeat chan<- struct{},
	r time.Time,
) {
	select {
	case <-done:
		return
	case <-pulse:
		sendPulse(heartbeat)
	case result <- r:
	}
}

func dowork(done <-chan struct{}) (<-chan struct{}, <-chan int) {
	heartbeat := make(chan struct{}, 1)
	workStream := make(chan int)

	go func() {
		defer close(heartbeat)
		defer close(workStream)

		for i := 0; i < 10; i++ {
			select {
			case heartbeat <- struct{}{}:
			default:
			}

			select {
			case <-done:
				return
			case workStream <- rand.Intn(10):
			}
		}
	}()

	return heartbeat, workStream
}

func main() {
	done := make(chan struct{})
	defer close(done)

	heartbeat, results := dowork(done)
	for {
		select {
		case _, ok := <-heartbeat:
			if !ok {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if !ok {
				return
			}
			fmt.Printf("results %v\n", r)
		}
	}
}
