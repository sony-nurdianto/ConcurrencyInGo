package main

import "fmt"

func main() {
	data := make(chan []int, 4)

	data <- []int{1, 2, 3, 4, 5}
	close(data)

	loopData := func(handleData chan<- int) {
		defer close(handleData)
		for i := range <-data {
			handleData <- i
		}
	}

	handleData := make(chan int)
	go loopData(handleData)

	for num := range handleData {
		fmt.Println(num)
	}
}
