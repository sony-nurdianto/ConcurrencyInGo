package main

import "fmt"

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

func repeat(done <-chan struct{}, values ...interface{}) <-chan interface{} {
	valueStream := make(chan interface{})

	go func() {
		defer close(valueStream)
		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case valueStream <- v:
				}
			}
		}
	}()
	return valueStream
}

func main() {
	done := make(chan struct{})
	defer close(done)

	var message string
	for token := range toString(done, take(done, repeat(done, "I", "Am", "Human"), 5)) {
		message += token
	}

	fmt.Printf("message: %s\n", message)
}
