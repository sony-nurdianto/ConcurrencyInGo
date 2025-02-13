package main

import (
	"fmt"
	"sync"
)

func main() {
	salutations := []string{"hello", "greetings", "good days"}

	wg := Salute(salutations)
	wg.Wait()
}

func Salute(s []string) *sync.WaitGroup {
	var wg sync.WaitGroup

	wg.Add(len(s))

	for _, salut := range s {
		go func(args string) {
			defer wg.Done()
			fmt.Println(args)
		}(salut)
	}

	return &wg
}
