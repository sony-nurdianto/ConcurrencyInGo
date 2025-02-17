package main

func orDone(done <-chan struct{}, c <-chan interface{}) <-chan interface{} {
	valueStream := make(chan interface{})

	go func() {
		defer close(valueStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if !ok {
					return
				}
				select {
				case valueStream <- v:
				case <-done:
				}
			}
		}
	}()

	return valueStream
}

func tee(done <-chan struct{}, in <-chan interface{}) (_, _ <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})

	go func() {
		defer close(out1)
		defer close(out2)

		for val := range orDone(done, in) {
			ch1, ch2 := out1, out2
			for i := 0; i < 2; i++ {
				select {
				case <-done:
				case ch1 <- val:
					ch1 = nil
				case ch2 <- val:
					ch2 = nil
				}
			}
		}
	}()

	return out1, out2
}
