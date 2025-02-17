package main

import (
	"fmt"
	"math/rand"
)

func repeatfn(done <-chan struct{}, fn func() interface{}) <-chan interface{} {
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

func toString(done <-chan struct{}, valueStream <-chan interface{}) <-chan string {
	stringStream := make(chan string)
	go func() {
		defer close(stringStream)
		for v := range valueStream {
			select {
			case <-done:
				return
			case stringStream <- v.(string):
			}
		}
	}()

	return stringStream
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

func main() {
	done := make(chan struct{})
	defer close(done)

	rand := func() interface{} {
		return rand.Int()
	}

	for num := range take(done, repeatfn(done, rand), 10) {
		fmt.Println(num)
	}
}
