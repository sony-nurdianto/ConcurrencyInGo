package main

import (
	"fmt"
	"net/http"
)

type Result struct {
	Error    error
	Response *http.Response
}

var checkStatus = func(done <-chan struct{}, urls ...string) <-chan Result {
	results := make(chan Result)
	go func() {
		defer close(results)
		for _, url := range urls {
			var result Result
			resp, err := http.Get(url)
			result = Result{Error: err, Response: resp}
			select {
			case <-done:
				return
			case results <- result:
			}
		}
	}()
	return results
}

func main() {
	done := make(chan struct{})
	defer close(done)

	errCount := 0

	urls := []string{"https://www.google.com", "https://badhost", "c", "d"}
	for response := range checkStatus(done, urls...) {
		if response.Error != nil {
			fmt.Printf("error: %v\n", response.Error)
			errCount++
			if errCount > 3 {
				fmt.Println("Too many errors, breaking!")
				break
			}
			continue
		}
		fmt.Printf("Response: %s\n", response.Response.Status)
	}
}
