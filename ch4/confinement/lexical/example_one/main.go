package main

import "fmt"

func main() {
	chanOwner := func() <-chan int {
		result := make(chan int, 5)
		go func() {
			defer close(result)
			for i := 0; i <= 5; i++ {
				result <- i
			}
		}()

		return result
	}

	consumer := func(result <-chan int) {
		for result := range result {
			fmt.Printf("Received: %d\n", result)
		}
		fmt.Println("Done Receiving")
	}

	result := chanOwner()
	consumer(result)
}
