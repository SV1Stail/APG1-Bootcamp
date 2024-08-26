package heap

import (
	"container/heap"
	"fmt"
)

type Present struct {
	Value int
	Size  int
}
type PresentHeap []Present

func (h PresentHeap) Len() int {
	return len(h)
}

func (h PresentHeap) Less(i, j int) bool {
	if h[i].Value > h[j].Value {
		return true
	} else if h[i].Value == h[j].Value {
		return h[i].Size > h[j].Size
	}
	return false
}

func (h PresentHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *PresentHeap) Push(x any) {
	*h = append(*h, x.(Present))
}

func (h *PresentHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

func GetNCoolestPresents(presents []Present, N int) ([]Present, error) {
	if N < 0 || N > len(presents) {
		return []Present{}, fmt.Errorf("Bad N")
	}
	ph := &PresentHeap{}
	heap.Init(ph)
	for _, present := range presents {
		heap.Push(ph, present)
	}

	var presentsRes []Present

	for i := 0; i < N; i++ {
		presentsRes = append(presentsRes, heap.Pop(ph).(Present))
	}
	return presentsRes, nil
}
