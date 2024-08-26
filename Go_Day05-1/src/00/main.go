package main

import (
	. "00/tree"
	"fmt"
)

func main() {
	tree := &TreeNode{
		Value: 1,
		Left: &TreeNode{
			Value: 0,
			Left: &TreeNode{
				Value: 0,
				Left:  nil,
				Right: nil,
			},
			Right: &TreeNode{Value: 1, Left: nil, Right: nil},
		},
		Right: &TreeNode{
			Value: 1,
			Left:  nil,
			Right: nil,
		},
	}
	if !AreToysBalanced(tree) {
		fmt.Println("Expected true, but got false")
	} else {
		fmt.Println("ok")
	}
}
