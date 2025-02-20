package main

import (
	"testing"
	"time"
)

func TestDoWorkGeneratesAllNumbers(t *testing.T) {
	done := make(chan struct{})
	defer close(done)

	intSlice := []int{1, 2, 3, 5}
	const timeout = 2 * time.Second

	heartbeat, results := DoWork(done, timeout/2, intSlice...)

	<-heartbeat

	i := 0

	for {
		select {
		case r, ok := <-results:
			if ok == false {
				return
			} else if expected := intSlice[i]; r != expected {
				t.Errorf("index %v: expected %v, but received %v,", i, expected, r)
			}
			i++
		case <-heartbeat:
		case <-time.After(timeout):
			t.Fatalf("test timeout")
		}
	}
}
