package firstfunc

import "sort"

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
