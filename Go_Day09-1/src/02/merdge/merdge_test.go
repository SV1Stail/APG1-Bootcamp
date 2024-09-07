package merdge_test

import (
	"02/merdge"
	"sync"
	"testing"
)

func TestMerdge1(t *testing.T) {
	ch1 := make(chan interface{})
	ch2 := make(chan interface{})
	ch3 := make(chan interface{})
	var wg sync.WaitGroup
	var mu sync.Mutex
	var mas [30]int
	index := 0
	insert := func(ch1 chan interface{}, s int, l int) {
		defer close(ch1)
		defer wg.Done()
		for i := s; i < l; i++ {
			ch1 <- i
			mu.Lock()
			mas[index] = i
			index++
			mu.Unlock()
		}
	}
	for i := 0; i < 30; i += 10 {
		wg.Add(1)
		if i < 10 {
			go insert(ch1, i, i+10)
		} else if i < 20 {
			go insert(ch2, i, i+10)
		} else {
			go insert(ch3, i, i+10)
		}
	}

	res := merdge.Multiplex(ch3, ch2, ch1)
	go func() {
		wg.Wait()
	}()
	resultSet := make(map[int]bool)
	for r := range res {
		if val, ok := r.(int); ok {
			resultSet[val] = true
		}
	}

	if len(resultSet) != 30 {
		t.Error("ERROR")
	}
	for _, val := range mas {
		if !resultSet[val] {
			t.Error("ERROR  2")
		}
	}
}
