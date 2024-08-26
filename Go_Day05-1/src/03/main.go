package main

import (
	"fmt"
	"os"
)

type Present struct {
	Value int
	Size  int
}

func grabPresents(presents []Present, diskSize int) ([]Present, error) {
	if diskSize < 1 {
		return []Present{}, fmt.Errorf("wrong disk size")
	}
	n := len(presents)
	if n < 1 {
		return []Present{}, fmt.Errorf("no Presents")
	}
	var mat = make([][]int, n+1)
	fmt.Println(mat)
	for i := range mat {
		mat[i] = make([]int, diskSize+1)
	}

	for i := 1; i <= n; i++ {
		for w := 1; w <= diskSize; w++ {
			if presents[i-1].Size <= w {
				mat[i][w] = max(mat[i-1][w], mat[i-1][w-presents[i-1].Size]+presents[i-1].Value)
			} else {
				mat[i][w] = mat[i-1][w]
			}
		}
	}
	result := []Present{}
	w := diskSize
	for i := n; i > 0 && w > 0; i-- {
		if mat[i][w] != mat[i-1][w] {
			result = append(result, presents[i-1])
			w -= presents[i-1].Size
		}
	}
	return result, nil
}

func main() {
	mat := []Present{
		{Value: 4, Size: 5},
		{Value: 3, Size: 1},
		{Value: 5, Size: 2},
		{Value: 5, Size: 1},
	}
	mat2, err := grabPresents(mat, 2)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(mat2)

}
