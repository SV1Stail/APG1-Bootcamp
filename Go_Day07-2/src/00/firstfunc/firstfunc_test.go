package firstfunc_test

import (
	"00/firstfunc"
	"reflect"
	"testing"
)

// func TestMinCoins1(t *testing.T) {
// 	tests := []struct {
// 		amount   int
// 		coins    []int
// 		expected []int
// 	}{
// 		{13, []int{1, 5, 10}, []int{10, 1, 1, 1}},
// 		{13, []int{1, 3, 4}, []int{4, 4, 4, 1}},
// 		{7, []int{1, 3, 4}, []int{4, 3}},
// 		{18, []int{1, 5, 10}, []int{10, 5, 1, 1, 1}},
// 		{23, []int{1, 7, 10}, []int{10, 10, 1, 1, 1}},
// 		{100, []int{1, 5, 10, 25, 100}, []int{100}},
// 		{15, []int{2, 3, 6, 7}, []int{7, 7, 2}}, // Проверка на возможность собрать с минимальным количеством элементов
// 		{9, []int{2, 3, 5}, []int{5, 3, 2}},
// 		{1, []int{2, 5, 10}, []int{}}, // Невозможно собрать сумму
// 		{100, []int{1, 5, 20, 50}, []int{50, 50}},
// 	}
// 	for _, test := range tests {
// 		result := firstfunc.MinCoins(test.amount, test.coins)
// 		if !reflect.DeepEqual(result, test.expected) {
// 			t.Errorf("For amount %d and coins %v, expected %v, but got %v", test.amount, test.coins, test.expected, result)
// 		}
// 	}

// }
func TestMinCoins2(t *testing.T) {
	tests := []struct {
		amount   int
		coins    []int
		expected []int
	}{
		{13, []int{1, 5, 10}, []int{10, 1, 1, 1}},
		{13, []int{1, 3, 4}, []int{4, 4, 4, 1}},
		{7, []int{1, 3, 4}, []int{4, 3}},
		{18, []int{1, 5, 10}, []int{10, 5, 1, 1, 1}},
		{23, []int{1, 7, 10}, []int{10, 10, 1, 1, 1}},
		{100, []int{1, 5, 10, 25, 100}, []int{100}},
		{15, []int{2, 3, 6, 7}, []int{7, 7, 2}}, // Проверка на возможность собрать с минимальным количеством элементов
		{9, []int{2, 3, 5}, []int{5, 3, 2}},
		{1, []int{2, 5, 10}, []int{2}}, // Невозможно собрать сумму
		{100, []int{1, 5, 20, 50}, []int{50, 50}},
	}
	for _, test := range tests {
		result := firstfunc.MinCoins2(test.amount, test.coins)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("For amount %d and coins %v, expected %v, but got %v", test.amount, test.coins, test.expected, result)
		}
	}

}
