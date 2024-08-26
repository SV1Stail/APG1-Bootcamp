package main

import (
	. "01/tree"
	"fmt"
	"os"
)

func main() {
	tree := &TreeNode{
		Value: true,
		Left: &TreeNode{
			Value: true,
			Left: &TreeNode{
				Value: true,
				Left:  nil,
				Right: nil,
			},
			Right: &TreeNode{Value: false, Left: nil, Right: nil},
		},
		Right: &TreeNode{
			Value: false,
			Left:  &TreeNode{Value: true},
			Right: &TreeNode{Value: true},
		},
	}

	mas := []bool{true, true, false, true, true, false, true}
	mas2 := UnrollGarland(tree)

	if len(mas) == len(mas2) {
		for i := 0; i < len(mas); i++ {
			if mas[i] != mas2[i] {
				fmt.Println("error", i)
				os.Exit(1)
			}
		}
	}
	fmt.Println("ok")
}
