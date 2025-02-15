package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		ch <- "pesan diterima"
	}()

	select {
	case msg := <-ch:
		fmt.Println(msg)
		return
	default:
		fmt.Println("Tidak ada pesan yang diterima")
	}
}
