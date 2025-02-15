package main

import (
	"fmt"
	"time"
)

func task(name string, ch chan<- string) {
	time.Sleep(time.Second)
	ch <- name + " selesai"
}

func main() {
	chCap := 2
	ch1 := make(chan string)
	ch2 := make(chan string)

	go task("Task 1", ch1)
	go task("Task 2", ch2)

	for i := 0; i < chCap; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println(msg1)
		case msg2 := <-ch2:
			fmt.Println(msg2)
		}
	}
}
