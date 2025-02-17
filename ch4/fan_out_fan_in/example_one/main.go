package main

import (
	"fmt"
	"sync"
)

func repeatFn(done <-chan struct{}, fn func() int) <-chan int {
	valueStream := make(chan int)

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

func take(done <-chan struct{}, valueStream <-chan int, workernums int) <-chan int {
	takeStream := make(chan int)

	go func() {
		defer close(takeStream)
		for i := 0; i < workernums; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()

	return takeStream
}

func generator(done <-chan struct{}, n ...int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for _, num := range n {
			select {
			case <-done:
				return
			case out <- num:
			}
		}
	}()
	return out
}

func multiplex(done <-chan struct{}, wg *sync.WaitGroup, multiplexedStream chan<- int, c <-chan int) {
	defer wg.Done()

	for i := range c {
		select {
		case <-done:
			return
		case multiplexedStream <- i:
		}
	}
}

func merge(done <-chan struct{}, channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	multiplexedStream := make(chan int)

	wg.Add(len(channels))

	for _, c := range channels {
		go multiplex(done, &wg, multiplexedStream, c)
	}

	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

func squearedNum(done <-chan struct{}, stream <-chan int) <-chan int {
	valueStream := make(chan int)

	go func() {
		defer close(valueStream)

		for n := range stream {
			select {
			case <-done:
				return
			case valueStream <- n * n:
			}
		}
	}()

	return valueStream
}

func main() {
	done := make(chan struct{})
	defer close(done)

	numGen := generator(done, 1, 2, 3, 4, 5, 6, 7, 9)

	for v := range take(done, squearedNum(done, numGen), 3) {
		fmt.Println(v)
	}
}
