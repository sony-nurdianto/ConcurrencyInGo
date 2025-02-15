package main

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
)

func UnblockingConnHelloWorld() {
	var wg sync.WaitGroup
	helloWorld := make(chan string)

	words := strings.Split("Hello World", " ")

	wg.Add(len(words))

	for _, v := range words {
		go func(w string) {
			defer wg.Done()
			helloWorld <- w
		}(v)
	}

	go func() {
		wg.Wait()
		close(helloWorld)
	}()

	var buf bytes.Buffer

	for value := range helloWorld {
		fmt.Fprintf(&buf, "%s ", value)
	}
}

func ConHelloWorld() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		ch1 <- "Hallo"
	}()

	go func() {
		hello := <-ch1
		ch2 <- fmt.Sprintf("%s World !!!", hello)
	}()

	<-ch2
}

func HelloWorld() {
	hello := "Hello"
	world := "World"

	_ = fmt.Sprintf("%s %s", hello, world)
}

func main() {
	var wg sync.WaitGroup
	helloWorld := make(chan string, 2)

	words := strings.Split("Hello World", " ")

	wg.Add(len(words))

	for _, v := range words {
		go func(w string) {
			defer wg.Done()
			helloWorld <- w
		}(v)
	}

	go func() {
		wg.Wait()
		close(helloWorld)
	}()

	var buf bytes.Buffer

	for value := range helloWorld {
		fmt.Fprintf(&buf, "%s ", value)
	}
	fmt.Println(buf.String())
}
