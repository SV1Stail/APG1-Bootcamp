package getelem

import (
	"fmt"
	"iter"
)

// (int, error) {
func GetElement(arr []int, idx int) (int, error) {
	if len(arr) <= idx {
		return 0, fmt.Errorf("ERROR: len arr < idx")
	} else if len(arr) == 0 {
		return 0, fmt.Errorf("lERROR: en(arr) == 0")

	} else if idx < 0 {
		return 0, fmt.Errorf("ERROR: idx < 0")
	}

	seq := func(yield func(int) bool) {
		for _, v := range arr {
			if !yield(v) {
				return
			}
		}
	}
	res := -1
	next, stop := iter.Pull(seq)
	defer stop()
	for i := 0; i <= idx; i++ {
		res, _ = next()

	}
	return res, nil

}
