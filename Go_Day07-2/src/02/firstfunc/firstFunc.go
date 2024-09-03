// package firstfunc provides funcs for use a minimal amount of coins to
// avoid slowing everyone else down
package firstfunc

import (
	"sort"
)

// Defoult solution
func MinCoins(val int, coins []int) []int {
	res := make([]int, 0)
	i := len(coins) - 1
	for i >= 0 {
		for val >= coins[i] {
			val -= coins[i]
			res = append(res, coins[i])
		}
		i -= 1
	}
	return res
}

// New solution has check for lenght of input data,
// another cicle,
// if we don't have coin smaller than val, we'll take the smallest one from the “coins” piece
func MinCoins2(val int, coins []int) []int {
	if len(coins) < 1 || val < 1 {
		return []int{}
	}
	valTmp := val
	res := make([]int, 0)
	sort.Ints(coins)
	n := len(coins)
	for i := n - 1; i >= 0; i-- {
		for valTmp >= coins[i] {
			valTmp -= coins[i]
			res = append(res, coins[i])
		}
	}
	if valTmp != 0 {
		res = append(res, coins[0])
	}
	return res
}

// Used go doc for html generatio:
// из ~/02:
// $go doc -all ./firstfunc/ > doc.txt
// pandoc docs.txt -o new.html
// !!!SUCCES!!!
func Zero(x int) int {
	return x + 1
}
