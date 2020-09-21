package main

type ListNode struct {
	Val int
	Next *ListNode
}

type TreeNode struct {
	Val int
	Left *TreeNode
	Right *TreeNode
}

type Node struct {
	Val int
	Neighbors []*Node
}

func findMinimum(a int, b int) int {
	if b < a {
		return b
	}
	return a
}

func isSliceEq(a, b []string) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// 生成全零的二维数组
func generateTwoDimensionalSlice(row int, column int) [][]int {
	result := [][]int{}
	rowSlice := []int{}
	for j := 0; j < column; j++ {
		rowSlice = append(rowSlice, 0)
	}
	for i := 0; i < row; i++ {
		result = append(result, rowSlice)
	}

	return result
}