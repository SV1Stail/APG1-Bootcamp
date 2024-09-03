package main

import (
	"00/firstfunc"
	"fmt"
)

func main() {
	tests := struct {
		amount   int
		coins    []int
		expected []int
	}{
		13, []int{1, 5, 10}, []int{10, 1, 1, 1},
	}
	x := firstfunc.MinCoins2(tests.amount, tests.coins)
	fmt.Println(x)
}
