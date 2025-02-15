package main

import "fmt"

func chanOwner() <-chan int {
	resultStream := make(chan int, 5)
	go func() {
		defer close(resultStream)
		for i := 0; i < 5; i++ {
			resultStream <- i
		}
	}()

	return resultStream
}

func main() {
	resultStream := chanOwner()

	for result := range resultStream {
		fmt.Println(result)
	}
	fmt.Println("Done Receiving ...")
}
