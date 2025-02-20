package main

import "time"

func DoWork(done <-chan struct{}, pulseInterval time.Duration, nums ...int) (<-chan struct{}, <-chan int) {
	heartBeat := make(chan struct{}, 1)
	intStream := make(chan int)

	go func() {
		defer close(heartBeat)
		defer close(intStream)

		time.Sleep(2 * time.Second)

		pulse := time.Tick(pulseInterval)

	numloop:
		for _, n := range nums {
			for {
				select {
				case <-done:
					return
				case <-pulse:
					select {
					case heartBeat <- struct{}{}:
					default:
					}
				case intStream <- n:
					continue numloop
				}
			}
		}
	}()

	return heartBeat, intStream
}

func main() {
}
