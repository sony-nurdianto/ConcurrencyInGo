package main

import (
	"fmt"
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

func dowork(done <-chan struct{}, pulseInterval time.Duration) (<-chan struct{}, <-chan time.Time) {
	heartbeat := make(chan struct{})
	results := make(chan time.Time)

	go func() {
		// defer close(heartbeat)
		// defer close(results)

		pulse := time.Tick(pulseInterval)
		workGen := time.Tick(2 * pulseInterval)

		for i := 0; i < 2; i++ {
			select {
			case <-done:
				return
			case <-pulse:
				sendPulse(heartbeat)
			case r := <-workGen:
				go sendResults(done, results, pulse, heartbeat, r)
			}
		}
	}()

	return heartbeat, results
}

func main() {
	done := make(chan struct{})
	time.AfterFunc(10*time.Second, func() { close(done) })

	const timeout = 2 * time.Second

	heartbeat, results := dowork(done, timeout/2)
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
		case <-time.After(timeout):
			fmt.Println("worker goroutine is not healthy")
			return
		}
	}
}
