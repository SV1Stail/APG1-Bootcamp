package main

import (
	"00/sleepsort"
	"fmt"
)

func main() {
	slice := []int{9, 8, 7, 4, 5, 6, 3, 2, 1}
	ch := sleepsort.SleepSort(slice)
	for x := range ch {

		fmt.Println(x)
	}
}
