package merdge

import "sync"

func Multiplex(chans ...<-chan interface{}) <-chan interface{} {
	resChan := make(chan interface{})
	var wg sync.WaitGroup
	for _, ch := range chans {
		wg.Add(1)
		go func(c interface{}) {
			defer wg.Done()
			for val := range ch {
				resChan <- val
			}
		}(ch)
	}

	go func() {

		wg.Wait()

		defer close(resChan)
	}()
	return resChan
	
}
