package main

import (
	"log"
	"os"
	"time"
)

type StartGoroutineFn func(done <-chan interface{}, pulseInterval time.Duration) (heartbeat <-chan interface{})

func or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	ordone := make(chan interface{})

	go func() {
		defer close(ordone)

		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-or(append(channels[3:], ordone)...):
			}
		}
	}()

	return ordone
}

func NewSteward(timeout time.Duration, startGoroutine StartGoroutineFn) StartGoroutineFn {
	return func(done <-chan interface{}, pulseInterval time.Duration) <-chan interface{} {
		heartbeat := make(chan interface{})
		go func() {
			defer close(heartbeat)

			var wardDone chan interface{}
			var wardHearbeat <-chan interface{}

			startWard := func() {
				wardDone = make(chan interface{})
				wardHearbeat = startGoroutine(or(wardDone, done), timeout/2)
			}

			startWard()
			pulse := time.Tick(pulseInterval)

		monitorLoop:
			for {
				timeoutSignal := time.After(timeout)
				for {
					select {
					case <-pulse:
						select {
						case heartbeat <- struct{}{}:
						default:
						}
					case <-wardHearbeat:
						continue monitorLoop
					case <-timeoutSignal:
						log.Println("steward: ward unhealthy; restarting")
						close(wardDone)
						startWard()
						continue monitorLoop
					case <-done:
						return
					}
				}
			}
		}()
		return heartbeat
	}
}

func DoWork(done <-chan interface{}, _ time.Duration) <-chan interface{} {
	log.Println("ward: hello, I'am Iresponsible")
	go func() {
		<-done
		log.Println("ward: I'am halting")
	}()
	return nil
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	doworkWithSteward := NewSteward(4*time.Second, DoWork)

	done := make(chan interface{})

	time.AfterFunc(9*time.Second, func() {
		log.Println("main: halting steward and ward")
		close(done)
	})

	for range doworkWithSteward(done, 4*time.Second) {
	}

	log.Println("Done")
}
