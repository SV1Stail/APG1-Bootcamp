package sleepsort

import (
	"sync"
	"time"
)

func SleepSort(slice []int) <-chan int {
	n := len(slice)
	ch := make(chan int)
	var wg sync.WaitGroup

		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(t int) {
				time.Sleep(time.Second * time.Duration(t))
				ch <- t
				wg.Done()
			}(i)
		}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}
