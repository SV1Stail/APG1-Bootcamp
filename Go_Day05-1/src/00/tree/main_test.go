package tree_test

import (
	. "00/tree"
	"testing"
)

func TestAreToysBalanced1(t *testing.T) {
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
		t.Errorf("Expected true, but got false")
	}

}
func TestAreToysBalanced2(t *testing.T) {
	tree := &TreeNode{
		Value: 1,
		Left: &TreeNode{
			Value: 1,
			Left: &TreeNode{
				Value: 0,
				Left:  nil,
				Right: nil,
			},
			Right: &TreeNode{Value: 1, Left: nil, Right: nil},
		},
		Right: &TreeNode{
			Value: 0,
			Left:  &TreeNode{Value: 1},
			Right: &TreeNode{Value: 1},
		},
	}
	if !AreToysBalanced(tree) {
		t.Errorf("Expected true, but got false")
	}

}

func TestAreToysBalanced5(t *testing.T) {
	tree := &TreeNode{
		Value: 1,
		Left: &TreeNode{
			Value: 1,
			Left: &TreeNode{
				Value: 0,
				Left:  &TreeNode{Value: 1},
				Right: nil,
			},
			Right: &TreeNode{Value: 1, Left: nil, Right: &TreeNode{Value: 1}},
		},
		Right: &TreeNode{
			Value: 0,
			Left:  &TreeNode{Value: 1, Left: &TreeNode{Value: 1}, Right: &TreeNode{Value: 1}},
			Right: &TreeNode{Value: 1},
		},
	}
	if !AreToysBalanced(tree) {
		t.Errorf("Expected true, but got false")
	}

}

func TestAreToysBalanced3(t *testing.T) {
	tree := &TreeNode{
		Value: 1,
		Left: &TreeNode{
			Value: 1,
		},
		Right: &TreeNode{
			Value: 0,
		},
	}
	if AreToysBalanced(tree) {
		t.Errorf("Expected false, but got true")
	}

}
func TestAreToysBalanced4(t *testing.T) {
	tree := &TreeNode{
		Value: 1,
		Left: &TreeNode{
			Value: 1,
			Right: &TreeNode{Value: 1},
		},
		Right: &TreeNode{
			Value: 0,
			Right: &TreeNode{Value: 1},
		},
	}
	if AreToysBalanced(tree) {
		t.Errorf("Expected false, but got true")
	}

}

func TestAreToysBalanced6(t *testing.T) {
	tree := &TreeNode{
		Value: 1,
		Left: &TreeNode{
			Value: 0,
			Left: &TreeNode{
				Value: 1,
				Left:  &TreeNode{Value: 1},
				Right: nil,
			},
			Right: &TreeNode{Value: 1, Left: nil, Right: &TreeNode{Value: 1}},
		},
		Right: &TreeNode{
			Value: 0,
			Left:  &TreeNode{Value: 1, Left: &TreeNode{Value: 1}, Right: &TreeNode{Value: 0}},
			Right: &TreeNode{Value: 1},
		},
	}
	if AreToysBalanced(tree) {
		t.Errorf("Expected false, but got true")
	}

}
