package tree

type TreeNode struct {
	Value int
	Left  *TreeNode
	Right *TreeNode
}

func countToys(tree *TreeNode) int {
	if tree == nil {
		return 0
	}
	return tree.Value + countToys(tree.Right) + countToys(tree.Left)
}

func AreToysBalanced(tree *TreeNode) bool {

	if tree == nil {
		return true
	}
	return countToys(tree.Left) == countToys(tree.Right)
}

// func main() {

// 	//     0
// 	//    / \
// 	//   0   1
// 	//  / \
// 	// 0   1
// 	tree := &TreeNode{
// 		Value: 1,
// 		Left: &TreeNode{
// 			Value: 0,
// 			Left: &TreeNode{
// 				Value: 0,
// 				Left:  nil,
// 				Right: nil,
// 			},
// 			Right: &TreeNode{Value: 1, Left: nil, Right: nil},
// 		},
// 		Right: &TreeNode{
// 			Value: 1,
// 			Left:  nil,
// 			Right: nil,
// 		},
// 	}
// 	fmt.Println(AreToysBalanced(tree))
// }
