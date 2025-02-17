package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func repeatFn(done <-chan struct{}, fn func() interface{}) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)

		for {
			select {
			case <-done:
				return
			case valueStream <- fn():
			}
		}
	}()

	return valueStream
}

func take(done <-chan struct{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{})

	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()

	return takeStream
}

func toInt(done <-chan struct{}, valueStream <-chan interface{}) <-chan int {
	intStream := make(chan int)

	go func() {
		defer close(intStream)
		for v := range valueStream {
			select {
			case <-done:
				return
			case intStream <- v.(int):
			}
		}
	}()

	return intStream
}

func primeFinder(done <-chan struct{}, intStream <-chan int) <-chan interface{} {
	primeStream := make(chan interface{})

	go func() {
		for integer := range intStream {
			integer -= 1
			prime := true

			for divisor := integer - 1; divisor > 1; divisor-- {
				if integer%divisor == 0 {
					prime = false
					break
				}
			}

			if prime {
				select {
				case <-done:
					return
				case primeStream <- integer:
				}
			}
		}
	}()

	return primeStream
}

func multiplexed(wg *sync.WaitGroup, done <-chan struct{}, multiplexedStream chan<- interface{}, c <-chan interface{}) {
	defer wg.Done()

	for i := range c {
		select {
		case <-done:
			return
		case multiplexedStream <- i:
		}
	}
}

func fanIn(done <-chan struct{}, channels ...<-chan interface{}) <-chan interface{} {
	var wg sync.WaitGroup
	multiplexedStream := make(chan interface{})

	wg.Add(len(channels))

	for _, c := range channels {
		go multiplexed(&wg, done, multiplexedStream, c)
	}

	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

func main() {
	done := make(chan struct{})
	defer close(done)

	startTime := time.Now()

	rand := func() interface{} {
		return rand.Intn(50_000_000)
	}

	randIntStream := toInt(done, repeatFn(done, rand))

	numFinder := runtime.NumCPU()
	fmt.Printf("Spining Up %d prime finders.\n", numFinder)
	finders := make([]<-chan interface{}, numFinder)
	fmt.Println("Primes:")
	for i := 0; i < numFinder; i++ {
		finders[i] = primeFinder(done, randIntStream)
	}
	for prime := range take(done, fanIn(done, finders...), 10) {
		fmt.Printf("\t%d\n", prime)
	}
	fmt.Printf("Search took: %v\n", time.Since(startTime))
}
