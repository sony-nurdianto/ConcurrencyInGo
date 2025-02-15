package main

import (
	"fmt"
	"time"
)

func worker(stopChan <-chan struct{}) {
	for {
		select {
		case <-stopChan:
			fmt.Println("Worker cancellation")
			return
		default:
			fmt.Println("Worker sedang Bekerja")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	stopChan := make(chan struct{})
	go worker(stopChan)
	stopChan <- struct{}{}
	time.Sleep(1 * time.Second)
}
