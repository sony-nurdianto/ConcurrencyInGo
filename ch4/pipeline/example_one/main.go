package main

import "fmt"

func generator(done <-chan struct{}, integers ...int) <-chan int {
	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for _, i := range integers {
			select {
			case <-done:
				return
			case intStream <- i:
			}
		}
	}()

	return intStream
}

func multiply(done <-chan struct{}, intstream <-chan int, multiplier int) <-chan int {
	multipliedStream := make(chan int)
	go func() {
		defer close(multipliedStream)
		for i := range intstream {
			select {
			case <-done:
				return
			case multipliedStream <- i * multiplier:
			}
		}
	}()
	return multipliedStream
}

func add(done <-chan struct{}, intStream <-chan int, additive int) <-chan int {
	additiveStream := make(chan int)

	go func() {
		defer close(additiveStream)
		for i := range intStream {
			select {
			case <-done:
				return
			case additiveStream <- i + additive:
			}
		}
	}()

	return additiveStream
}

func main() {
	done := make(chan struct{})
	defer close(done)

	intStream := generator(done, 1, 2, 3, 4, 5, 6)
	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)
	for v := range pipeline {
		fmt.Println(v)
	}
}
