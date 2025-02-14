package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	mutex sync.Mutex
	cond  = sync.NewCond(&mutex)
	ready = false
)

func producer() {
	mutex.Lock()
	fmt.Println("Producing data ...")
	time.Sleep(2 * time.Second) // simulation production process
	ready = true
	mutex.Unlock()
	cond.Broadcast()
}

func consumer(id int) {
	mutex.Lock()
	for !ready {
		fmt.Printf("Consumer %d waiting...", id)
		cond.Wait()
	}
	fmt.Printf("Consumer %d data!\n", id)
	mutex.Unlock()
}

func main() {
	var wg sync.WaitGroup

	waitProcess := 3

	wg.Add(3)

	for i := 0; i < waitProcess-1; i++ {
		go func() {
			defer wg.Done()
			consumer(i)
		}()
	}

	go func() {
		defer wg.Done()
		producer()
	}()

	wg.Wait()
}
