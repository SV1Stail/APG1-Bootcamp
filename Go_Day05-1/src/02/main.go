package main

import (
	h "02/heap"
	"fmt"
	"os"
)

func main() {
	presents := []h.Present{
		{Value: 5, Size: 1},
		{Value: 4, Size: 5},
		{Value: 3, Size: 1},
		{Value: 5, Size: 2},
	}

	coolest, err := h.GetNCoolestPresents(presents, 2)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, pr := range coolest {
		fmt.Printf("Value: %d, Size: %d\n", pr.Value, pr.Size)

	}
}
