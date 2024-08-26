package tree

type TreeNode struct {
	Value bool
	Left  *TreeNode
	Right *TreeNode
}

func UnrollGarland(tree *TreeNode) []bool {
	if tree == nil {
		return []bool{}
	}

	var result []bool
	queue := []*TreeNode{tree}
	level := 1
	for len(queue) > 0 {
		resLen := len(result)
		width := len(queue)
		for i := 0; i < width; i++ {
			node := queue[0]
			queue = queue[1:]
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
			result = append(result, node.Value)
		}
		if level%2 != 0 {
			for i, j := resLen, len(result)-1; i < j; i, j = i+1, j-1 {
				result[i], result[j] = result[j], result[i]
			}
		}
		level++
	}
	return result
}

// func UnrollGarland(tree *TreeNode) []bool {
// 	if tree == nil {
// 		return []bool{}
// 	}
// 	var result []bool
// 	queue := []*TreeNode{tree}
// 	level := 1
// 	for len(queue) > 0 {
// 		curWidth := len(queue)
// 		curLevelResult := []bool{}
// 		for i := 0; i < curWidth; i++ {
// 			node := queue[0]
// 			queue := queue[:1]
// 			if node.Left != nil {
// 				queue = append(queue, node.Left)
// 			}
// 			if node.Right != nil {
// 				queue = append(queue, node.Right)
// 			}
// 			curLevelResult = append(curLevelResult, node.Value)
// 		}
// 		if level%2 != 0 {
// 			for i, j := 0, len(curLevelResult); i < j; i, j = i+1, j-1 {
// 				curLevelResult[i], curLevelResult[j] = curLevelResult[j], curLevelResult[i]
// 			}
// 		}
// 		result = append(result, curLevelResult...)
// 		level++
// 	}
// 	return result
// }
