package main

import (
	"fmt"
	"math"
)

// 深拷贝无向图
func cloneGraph(node *Node) *Node {
	if node == nil {
		return nil
	}
	copyNode := new(Node)
	hasDfsNodeMap := make(map[int]*Node)
	hasUsedNodeMap := make(map[int]*Node)
	dfsNodeNeighbor(hasUsedNodeMap, hasDfsNodeMap, node, copyNode)
	return copyNode
}

func dfsNodeNeighbor(hasUsedNodeMap map[int]*Node, hasDfsNodeMap map[int]*Node, originNode *Node, copyNode *Node) {
	_, ok := hasDfsNodeMap[originNode.Val]
	if ok {
		return
	}
	copyNode.Val = originNode.Val
	hasDfsNodeMap[originNode.Val] = copyNode
	hasUsedNodeMap[originNode.Val] = copyNode
	for i := 0; i < len(originNode.Neighbors); i++ {
		neighborNode := originNode.Neighbors[i]
		neighborNodeCopy, ok := hasUsedNodeMap[neighborNode.Val]
		if !ok {
			neighborNodeCopy = new(Node)
			neighborNodeCopy.Val = neighborNode.Val
		}
		hasUsedNodeMap[neighborNode.Val] = neighborNodeCopy
		copyNode.Neighbors = append(copyNode.Neighbors, neighborNodeCopy)
		dfsNodeNeighbor(hasUsedNodeMap, hasDfsNodeMap, neighborNode, neighborNodeCopy)
	}
}

// 前序遍历二叉树
func preorderTraversal(root *TreeNode) []int {
	result := []int{}
	lookSubtree(&result, root)
	return result
}

func lookSubtree(result *[]int, node *TreeNode) {
	if node == nil {
		return
	}
	*result = append(*result, node.Val)
	if node.Left != nil {
		lookSubtree(result, node.Left)
	}
	if node.Right != nil {
		lookSubtree(result, node.Right)
	}
}

// 后序遍历二叉树
func postorderTraversal(root *TreeNode) []int {
	result := []int{}
	lookSubtree2(&result, root)
	return result
}

func lookSubtree2(result *[]int, node *TreeNode) {
	if node == nil {
		return
	}
	if node.Left != nil {
		lookSubtree(result, node.Left)
	}
	if node.Right != nil {
		lookSubtree(result, node.Right)
	}
	*result = append(*result, node.Val)
}

//LRU算法
type LinkNode struct {
	next *LinkNode
	last *LinkNode
	key int
	val int
}
type LRUCache struct {
	capacity int
	kvMap map[int]*LinkNode
	head *LinkNode
	end *LinkNode
	linkNodeSize int
}
func Constructor(capacity int) LRUCache {
	var obj LRUCache
	obj.capacity = capacity
	obj.kvMap = map[int]*LinkNode{}
	head := &LinkNode{}
	end := &LinkNode{}
	head.last = nil
	head.next = end
	end.last = head
	end.next = nil
	obj.head = head
	obj.end = end
	obj.linkNodeSize = 0
	return obj
}

func deleteNode(targetNode *LinkNode) {
	if targetNode.next != nil && targetNode.last != nil {
		targetNode.last.next = targetNode.next
		targetNode.next.last = targetNode.last
	}
	if targetNode.next == nil {
		targetNode.last.next = nil
	}
	if targetNode.last == nil {
		targetNode.next.last = nil
	}
	targetNode = nil
}

func deleteFirstNode(head *LinkNode) {
	head.next = head.next.next
	head.next.last = head
}
func insertToEnd(end *LinkNode, targetNode *LinkNode) {
	targetNode.next = end
	targetNode.last = end.last
	end.last.next = targetNode
	end.last = targetNode
}
func (this *LRUCache) Get(key int) int {
	node, ok := this.kvMap[key]
	if !ok {
		return -1
	}
	val := node.val
	node.last.next = node.next
	node.next.last = node.last
	this.end.last.next = node
	node.next = this.end
	node.last = this.end.last
	this.end.last = node
	return val
}
func (this *LRUCache) Put(key int, value int)  {
	targetNode, ok := this.kvMap[key]
	if this.linkNodeSize == this.capacity && !ok {
		delete(this.kvMap, this.head.next.key)
		deleteNode(this.head.next)
		this.linkNodeSize -= 1
	}
	if ok {
		deleteNode(targetNode)
		this.linkNodeSize -= 1
	}
	valueNode := &LinkNode{}
	valueNode.val = value
	valueNode.key = key
	this.kvMap[key] = valueNode
	insertToEnd(this.end, valueNode)
	this.linkNodeSize += 1
}


//type LRUCache struct {
//	capacity int
//	kvMap map[int]int
//	keyList []int
//}
//
//func Constructor(capacity int) LRUCache {
//	var obj LRUCache
//	obj.capacity = capacity
//	obj.kvMap = make(map[int]int)
//	obj.keyList = []int{}
//	return obj
//}
//
//
//func (this *LRUCache) Get(key int) int {
//	value, ok := this.kvMap[key]
//	if !ok {
//		return -1
//	}
//	updateList(this, key)
//	return value
//}
//
//func (this *LRUCache) Put(key int, value int)  {
//	_, ok := this.kvMap[key]
//	if len(this.kvMap) == this.capacity && !ok {
//		delete(this.kvMap, this.keyList[0])
//		this.keyList = this.keyList[1:]
//	}
//	this.kvMap[key] = value
//	updateList(this, key)
//}
//
//func updateList(this *LRUCache, key int) {
//	targetIndex := -1
//	listLength := len(this.keyList)
//	for i := 0; i < listLength; i++ {
//		if this.keyList[i] == key {
//			targetIndex = i
//			break
//		}
//	}
//	if targetIndex > -1 {
//		this.keyList = append(this.keyList[:targetIndex], this.keyList[targetIndex + 1:]...)
//	}
//	this.keyList = append(this.keyList, key)
//	fmt.Println(this.keyList)
//}

/**
 * Your LRUCache object will be instantiated and called as such:
 * obj := Constructor(capacity);
 * param_1 := obj.Get(key);
 * obj.Put(key,value);
 */

func pathSum(root *TreeNode, sum int) [][]int {
	paths := [][]int{}
	getPathSum(root, sum, []int{}, &paths)
	return paths
}

func getPathSum(root *TreeNode, sum int, path []int, paths *[][]int) bool {
	if root.Right == nil && root.Left == nil {
		if root.Val == sum {
			path = append(path, root.Val)
			*paths = append(*paths, path)
		}
		path = []int{}
		return root.Val == sum
	}

	if root.Left != nil {
		path = append(path, root.Val)
		_ = getPathSum(root.Left, sum - root.Val, path, paths)
	}

	if root.Right != nil {
		if root.Left == nil {
			path = append(path, root.Val)
		}
		_ = getPathSum(root.Right, sum - root.Val, path, paths)
	}
	path = []int{}
	return false
}

//
//给定一个二叉树和一个目标和，判断该树中是否存在根节点到叶子节点的路径，这条路径上所有节点值相加等于目标和。
//
//说明: 叶子节点是指没有子节点的节点。
//
//示例: 
//给定如下二叉树，以及目标和 sum = 22，
//
//	5
//	/ \
//	4   8
//	/   / \
//	11  13  4
//	/  \      \
//	7    2      1
//
//来源：力扣（LeetCode）
//链接：https://leetcode-cn.com/problems/path-sum
//
func hasPathSum(root *TreeNode, sum int) bool {
	if root.Right == nil && root.Left == nil {
		return root.Val == sum
	}
	result := false
	if root.Left != nil {
		result = hasPathSum(root.Left, sum - root.Val)
	}
	if result {
		return true
	}
	if root.Right != nil {
		return hasPathSum(root.Right, sum - root.Val)
	}
	return false
}

func minimumTotal1(triangle [][]int) int {
	//第一种写法
	rows := len(triangle)
	for i := rows - 1; i >= 0; i-- {
		for j := 0; j < i; j++ {
			triangle[i - 1][j] = findMinimum(triangle[i][j], triangle[i][j + 1]) + triangle[i - 1][j]
		}
	}
	return  triangle[0][0]
	//第二种写法
	//rows := len(triangle)
	//dp := make([]int, rows + 1)
	//for i := rows - 1; i >= 0; i-- {
	//	for j := 0; j < len(triangle[i]); j++ {
	//		dp[j] = findMinimum(dp[j], dp[j + 1]) + triangle[i][j]
	//	}
	//}
	//return dp[0]
}

// 找出三角形最小路径和 --- 超时
func minimumTotal(triangle [][]int) int {
	if len(triangle) == 0 {
		return 0
	}
	if len(triangle[0]) == 0 {
		return 0
	}

	pointer := &triangle
	return minimumTotalSplit(pointer, 0, 0)
}

func minimumTotalSplit(triangle *[][]int, line int, column int) int {
	if len(*triangle) == line + 1 {
		return (*triangle)[line][column]
	}
	total := (*triangle)[line][column]
	total += findMinimum(minimumTotalSplit(triangle, line + 1, column), minimumTotalSplit(triangle, line + 1, column + 1))
	return total
}

// 判断二叉树是否对称
func isSymmetric(root *TreeNode) bool {
	if root == nil {
		return true
	}

	return compareLeftAndRight(root.Left, root.Right)
}

func compareLeftAndRight(leftNode *TreeNode, rightNode *TreeNode) bool {
	if leftNode == nil && rightNode == nil {
		return true
	}
	if leftNode == nil || rightNode == nil || leftNode.Val != rightNode.Val {
		return false
	}
	return compareLeftAndRight(leftNode.Left ,rightNode.Right) && compareLeftAndRight(leftNode.Right, rightNode.Left)
}

// 层序遍历二叉树 - 优化 - BFS
func levelOrder1(root *TreeNode) [][]int {
	result := [][]int{}
	if root == nil {
		return result
	}
	nodes := []*TreeNode{root}
	for {
		if len(nodes) == 0 {
			break
		}
		sameLevelNodeNumber := len(nodes)
		sameLevelValList := []int{}
		for i := 0; i < sameLevelNodeNumber; i++ {
			firstNode := nodes[0]
			sameLevelValList = append(sameLevelValList, firstNode.Val)
			nodes = nodes[1:]
			if firstNode.Left != nil {
				nodes = append(nodes, firstNode.Left)
			}
			if firstNode.Right != nil {
				nodes = append(nodes, firstNode.Right)
			}
		}
		result = append(result, sameLevelValList)
	}
	return result
}

// 层序遍历二叉树
func levelOrder(root *TreeNode) [][]int {
	result := [][]int{}
	if root == nil {
		return result
	}
	sblingNode := []*TreeNode{root}
	valList := []int{}
	for {
		if len(sblingNode) == 0 {
			break
		}
		valList, sblingNode = getSblingValueList(sblingNode)
		result = append(result, valList)
	}
	return result
}

func getSblingValueList(nodes []*TreeNode) ([]int, []*TreeNode) {
	valList := []int{}
	nodeList := []*TreeNode{}
	for i := 0; i < len(nodes); i++ {
		node := nodes[i]
		valList = append(valList, node.Val)
		if node.Left != nil {
			nodeList = append(nodeList, node.Left)
		}
		if node.Right != nil {
			nodeList = append(nodeList, node.Right)
		}
	}
	return valList, nodeList
}

// 反转链表
func reverseList(head *ListNode) *ListNode {
	currentNode := &ListNode{}
	pointer := head
	for {
		if pointer.Next == nil {
			pointer.Next = currentNode
			head.Next = nil
			break
		}
		currentNode.Next = pointer
		currentNode = currentNode.Next
		pointer = pointer.Next
	}
	return currentNode
}

// 找出字符串的最长回文子串 - 优化 - 中心扩散
func longestPalindrome2(s string) string {
	if len(s) <= 1 {
		return s
	}
	maxCount := 1
	maxString := s[:1]
	for centerPointerIndex := 0; centerPointerIndex < len(s); centerPointerIndex++ {
		for moveIndex := 1; (centerPointerIndex + moveIndex) < len(s) && (centerPointerIndex - moveIndex) >= 0; moveIndex++ {
			if s[centerPointerIndex - moveIndex] != s[centerPointerIndex + moveIndex] {
				break
			}
			if moveIndex * 2 + 1 > maxCount {
				maxCount = moveIndex * 2 + 1
				maxString = s[centerPointerIndex - moveIndex: centerPointerIndex + moveIndex + 1]
			}
		}
	}
	for centerPointerIndex := 0.5; centerPointerIndex <= float64(len(s)) - 1.5; centerPointerIndex++ {
		for moveIndex := 0.5; (centerPointerIndex + moveIndex) <= float64(len(s) - 1) && (centerPointerIndex - moveIndex) >= 0; moveIndex++ {
			if s[int(centerPointerIndex - moveIndex)] != s[int(centerPointerIndex + moveIndex)] {
				break
			}
			if moveIndex * 2 + 1 > float64(maxCount) {
				maxCount = int(moveIndex * 2 + 1)
				maxString = s[int(centerPointerIndex - moveIndex): int(centerPointerIndex + moveIndex) + 1]
			}
		}
	}
	return maxString
}

// 找出字符串的最长回文子串
func longestPalindrome(s string) string {
	if len(s) == 0 {
		return s
	}
	maxCount := 1
	maxString := s[:1]
	for i := 0; i < len(s); i++ {
		for j := i + 1; j < len(s) + 1; j++ {
			currentString := s[i: j]
			fmt.Println(currentString)
			if !isPalindrome(currentString) || (j - i + 1) <= maxCount {
				continue
			}
			maxString = currentString
			maxCount = j - i + 1
		}
	}
	return maxString
}

// 判断一个数是否为快乐数
func isHappy(n int) bool {
	pointer0 := getNextValue(n)
	pointer1 := getNextValue(n)
	for {
		pointer0 = getNextValue(pointer0)
		pointer1 = getNextValue(getNextValue(pointer1))
		if pointer1 == 1 || pointer0 == 1 {
			return true
		}
		if pointer0 == pointer1 {
			return false
		}
	}
	return true
}

func getNextValue(n int) int {
	index := 1
	result := math.Pow(float64(n % 10), float64(2))
	for {
		if n % 10 == n {
			break
		}
		n = n / 10
		result += math.Pow(float64(n % 10), float64(2))
		index += 1
	}
	return int(result)
}

// 列出字符串所有回文串数组 - 优化3 - dfs回溯
func partition4(s string) [][]string {
	result := [][]string{}
	current := []string{}
	if len(s) == 0 {
		return result
	}
	dfs(s, 0, &result, current)
	return result
}

// 两个参数是否引用传参！！！
func dfs(s string, startIndex int, result *[][]string, current []string) {
	if startIndex == len(s) {
		// 一定要拷贝！！！！！
		var currentCopy = make([]string, len(current))
		copy(currentCopy, current)
		*result = append(*result, currentCopy)
		return
	}
	for i := startIndex; i < len(s); i++ {
		str1 := s[startIndex: i + 1]
		if isPalindrome(str1) {
			current = append(current, str1)
			dfs(s, i+1, result, current)
			current = current[:(len(current) - 1)]
		}
	}
}

// 列出字符串所有回文串数组 - 优化2 - 分治+动态规划
func partition3(s string) [][]string {
	return partitionHelper3(0, s, getDp(s))
}

func partitionHelper3(startIndex int, s string, dp [][]bool) [][]string {
	result := [][]string{}
	if startIndex == len(s) {
		return append(result, []string{})
	}
	for i := startIndex; i < len(s); i++ {
		if !dp[startIndex][i] {
			continue
		}
		restResult := partitionHelper3(i + 1, s, dp)
		for j := 0; j < len(restResult); j++ {
			result = append(result, append([]string{s[startIndex: i + 1]}, restResult[j]...))
		}
	}
	return result
}

func getDp(s string) [][]bool {
	slen := len(s)
	dp := make([][]bool, slen)
	for i := 0; i < slen; i++ {
		dp[i] = make([]bool, slen)
	}
	for i := 0; i < slen; i++ {
		for j := 0; j <= i; j++ {
			if i == j {
				dp[j][i] = true
			} else if i - j == 1 {
				dp[j][i] = s[i] == s[j]
			} else {
				dp[j][i] = s[i] == s[j] && dp[j+1][i-1]
			}
		}
	}
	return dp
}

// 列出字符串的所有回文串数组 - 优化1
func partition2(s string) [][]string {
	return partitionHelper2(0, s)
}

func partitionHelper2(startIndex int, s string) [][]string{
	result := [][]string{}
	if startIndex == len(s) {
		return append(result, []string{})
	}
	for i := startIndex; i < len(s); i++ {
		str1 := s[startIndex: i + 1]
		if isPalindrome(str1) {
			restResult := partitionHelper2(0, s[i + 1 : len(s)])
			for j := 0; j < len(restResult); j++ {
				result = append(result, append([]string{str1}, restResult[j]...))
			}
		}
	}
	return result
}


// 列出字符串的所有回文串数组
func partition(s string) [][]string {
	result := partitionHelper(s)
	newResult := [][]string{}
	for i := range result {
		element := result[i]
		if element[len(element) - 1] == "" {
			element = element[0: len(element) - 1]
		}
		hasRepeat := false
		for j := 0; j < len(newResult); j++ {
			if isSliceEq(newResult[j], element) {
				hasRepeat = true
				break
			}
		}
		if !hasRepeat {
			newResult = append(newResult, element)
		}
	}
	return newResult
}

func partitionHelper(s string) [][]string {
	result := [][]string{}
	if len(s) == 1 {
		return append(result, []string{s})
	}

	for i := 1; i <= len(s); i++ {
		str1 := s[0:i]
		if !isPalindrome(str1) {
			continue
		}
		if isPalindrome(s[i: len(s)]) {
			result = append(result, []string{str1, s[i: len(s)]})
		}
		restResult := partition(s[i: len(s)])
		for j := 0; j < len(restResult); j++ {
			tempResult := []string{str1}
			for k := 0; k < len(restResult[j]); k++ {
				if restResult[j][k] == " " {
					continue
				}
				tempResult = append(tempResult, restResult[j][k])
			}

			result = append(result, tempResult)
		}
	}
	return result
}

func isPalindrome(s string) bool {
	i := 0
	j := len(s) - 1
	result := true
	for {
		if  i >= j {
			break
		}
		if s[i] != s[j] {
			return false
		}
		i += 1
		j -= 1
	}
	return result
}

// 字符串字符是否可以组成回文字符串
func canPermutePalindrome(s string) bool {
	keymap := map[string]string{}
	for _, ch := range s {
		str := string(ch)
		val := keymap[str]
		if val == "" {
			keymap[str] = str
		} else {
			delete(keymap, str)
		}
	}
	return len(keymap) < 2
}


/*
给定一个范围在  1 ≤ a[i] ≤ n ( n = 数组大小 ) 的 整型数组，数组中的元素一些出现了两次，另一些只出现一次。

找到所有在 [1, n] 范围之间没有出现在数组中的数字。

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/find-all-numbers-disappeared-in-an-array
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
*/
func findDisappearedNumbers(nums []int) []int {
	var result []int
	for i := 0; i < len(nums); i++ {
		indexValue := nums[i]
		if indexValue < 0 {
			indexValue = -indexValue
		}
		targetIndexValue := nums[indexValue - 1]
		if targetIndexValue > 0 {
			targetIndexValue = -targetIndexValue
		}
		nums[indexValue - 1] = targetIndexValue
	}
	for i := 0; i < len(nums); i++ {
		if nums[i] >= 0 {
			result = append(result, i + 1)
		}
	}
	return result
}

// 删除倒数第n个节点
func removeNthFromEnd(head *ListNode, n int) *ListNode {
	pointer1 := head
	var pointer2 = &ListNode{}
	var pointer3 = &ListNode{}
	index := 1
	for {
		if pointer1.Next == nil {
			if index == n {
				return head.Next
			}
			pointer3.Next = pointer2.Next
			pointer2.Next = nil
			break
		}
		if index == n {
			pointer3 = head
		}
		if index == n - 1 {
			pointer2 = head
		}
		if index > n - 1 && pointer2.Next != nil {
			pointer2 = pointer2.Next
		}
		if index > n && pointer3.Next != nil {
			pointer3 = pointer3.Next
		}
		pointer1 = pointer1.Next
		index += 1
	}
	return head
}

// 盛水最多的容器 - 双指针
func maxArea2(height []int) int {
	var size = 0
	if len(height) < 2 {
		return size
	}
	var leftIndex = 0
	var rightIndex = len(height) - 1
	var lower = height[leftIndex]
	if lower > height[rightIndex] {
		lower = height[rightIndex]
	}
	size = lower * (rightIndex - leftIndex)
	for {
		if height[leftIndex] > height[rightIndex] {
			rightIndex -= 1
		} else {
			leftIndex += 1
		}
		if rightIndex == leftIndex {
			break
		}
		tempSize := (rightIndex - leftIndex) * height[rightIndex]
		if height[leftIndex] < height[rightIndex] {
			tempSize = (rightIndex - leftIndex) * height[leftIndex]
		}
		if tempSize > size {
			size = tempSize
		}
	}
	return size
}

// 盛水最多的容器 - 暴力破解
func maxArea(height []int) int {
	var size = 0
	for i := 0; i < len(height); i++ {
		for j := i+1; j < len(height); j++ {
			squareHeight := height[i]
			if height[i] > height[j] {
				squareHeight = height[j]
			}
			newSize := squareHeight * (j - i)
			if newSize > size {
				size = newSize
			}
		}
	}
	return size
}

// 找出只出现一次的数字 - 异或的特点
func singleNumber(nums []int) int {
	result := 0
	for index := 0; index < len(nums); index++ {
		result ^= nums[index]
	}
	return result
}

// 合并两个有序链表
func mergeTwoLists(l1 *ListNode, l2 *ListNode) *ListNode {
	nextNode := &ListNode{}
	pointer := nextNode

	for {
		if l1 == nil || l2 == nil {
			break
		}
		if l1.Val > l2.Val {
			pointer.Next = l2
			l2 = l2.Next
		} else {
			pointer.Next = l1
			l1 = l1.Next
		}

		pointer = pointer.Next
	}
	if l1 == nil && l2 != nil {
		pointer.Next = l2
	}
	if l2 == nil && l1 != nil {
		pointer.Next = l1
	}
	return nextNode.Next
}
