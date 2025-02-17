package main

import (
	"fmt"
	"sync"
)

// ✅ Perbaiki iterasi slice
func generator(done <-chan struct{}, n ...int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for _, num := range n { // ✅ Gunakan for _, num untuk iterasi slice
			select {
			case <-done:
				return
			case out <- num:
			}
		}
	}()
	return out
}

// ✅ Tambahkan sync.WaitGroup untuk menutup channel
func squareNum(done <-chan struct{}, numStream <-chan int, workerCount int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func() {
			defer wg.Done()
			for n := range numStream {
				select {
				case <-done:
					return
				case out <- n * n:
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out) // ✅ Tutup channel setelah semua worker selesai
	}()

	return out
}

// ✅ Perbaiki defer wg.Done() dan baca channel dalam loop
func merge(done <-chan struct{}, channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	wg.Add(len(channels))
	for _, ch := range channels {
		go func(c <-chan int) {
			defer wg.Done()
			for val := range c { // ✅ Iterasi agar tidak hanya membaca satu nilai
				select {
				case <-done:
					return
				case out <- val:
				}
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// ✅ Tambahkan fungsi tee untuk membagi channel
func tee(done <-chan struct{}, in <-chan int) (<-chan int, <-chan int) {
	out1 := make(chan int)
	out2 := make(chan int)

	go func() {
		defer close(out1)
		defer close(out2)
		for val := range in {
			select {
			case <-done:
				return
			case out1 <- val:
			}

			select {
			case <-done:
				return
			case out2 <- val:
			}
		}
	}()

	return out1, out2
}

func main() {
	done := make(chan struct{})
	defer close(done)

	numGen := generator(done, 1, 2, 3, 4, 5, 6, 7, 8, 9)

	// ✅ Gunakan tee untuk menduplikasi stream
	tee1, tee2 := tee(done, numGen)
	squared1 := squareNum(done, tee1, 3)
	squared2 := squareNum(done, tee2, 4)

	merg := merge(done, squared1, squared2)

	for m := range merg {
		fmt.Println(m)
	}
}
