package tree

import (
	"fmt"
	"math"
)

/**
 * B+ Tree Implementation
 * Reference: https://en.wikipedia.org/wiki/B%2B_tree
 * Following Sean Davis's Rules for B+ Tree Implementation, accessed via https://www.youtube.com/watch?v=49P_GDeMDRo
 *
**/

// BPNode represents a node in the B+ Tree
type BPNode struct {
	keys     []int
	children []*BPNode
	isLeaf   bool
	parent   *BPNode
	next     *BPNode // For leaf nodes to link to the next sibling
}

type BPTree struct {
	root *BPNode
	M    int // Maximum number of children in internal nodes, Minimum number of children in internal nodes is ceil(M/2)
	L    int // Maximum number of elements in leaf nodes, Minimum number of elements in leaf nodes is ceil(L/2)
}

type Tree interface {
	Insert(key int) error          // Returns an error if the key is already present
	Delete(key int) error          // Returns an error if the key is not found
	Search(key int) (*BPNode, int) // Returns the node and the index of the key in the node
	PrintTree()
	PrintTreeLevel(level int)
}

// CreateTree creates a new B+ Tree with the given maximum number of children in internal nodes (M) and maximum number of elements in leaf nodes (L)
func CreateTree(m, l int) *BPTree {
	return &BPTree{
		root: &BPNode{
			keys:   []int{},
			isLeaf: true,
		},
		M: m,
		L: l,
	}
}

func (tree *BPTree) MaxChildrenInternalNodes() int {
	return tree.M
}

func (tree *BPTree) MinChildrenInternalNodes() int {
	return int(math.Ceil(float64(tree.M) / 2))
}

func (tree *BPTree) MaxElementsLeafNodes() int {
	return tree.L
}

func (tree *BPTree) MinElementsLeafNodes() int {
	return int(math.Ceil(float64(tree.L) / 2))
}

func (tree *BPTree) Insert(value int) {
	fmt.Printf("Inserting value %d into the B+ tree\n", value)

	fmt.Println("Finding leaf node for insertion")
	node := findLeafNode(tree.root, value)
	fmt.Printf("Found leaf node with keys: %v\n", node.keys)

	fmt.Printf("Inserting value %d into leaf node\n", value)
	insertIntoLeafNode(node, value, tree.L)
	fmt.Printf("After insertion, leaf node keys: %v\n", node.keys)

	if len(node.keys) > tree.L {
		fmt.Printf("Leaf node overflow detected. Current keys: %v, Max allowed: %d\n", node.keys, tree.L)
		fmt.Println("Splitting leaf node")
		splitLeafNode(tree, node)
		fmt.Println("Leaf node split complete")
	} else {
		fmt.Println("No split required")
	}

	fmt.Printf("Insertion of value %d complete\n", value)
}

func findLeafNode(node *BPNode, value int) *BPNode {
	if node.isLeaf {
		return node
	}

	for i := 0; i < len(node.keys); i++ {
		if value < node.keys[i] {
			return findLeafNode(node.children[i], value)
		}
	}

	return findLeafNode(node.children[len(node.children)-1], value)
}

func insertIntoLeafNode(node *BPNode, value int, L int) {
	node.keys = append(node.keys, value)
	for i := len(node.keys) - 1; i > 0 && node.keys[i] < node.keys[i-1]; i-- {
		node.keys[i], node.keys[i-1] = node.keys[i-1], node.keys[i]
	}
}

func splitLeafNode(tree *BPTree, node *BPNode) {
	fmt.Println("Starting splitLeafNode function")
	fmt.Printf("Original node keys: %v\n", node.keys)

	newLeaf := &BPNode{
		keys:   make([]int, 0),
		isLeaf: true,
	}
	fmt.Println("Created new leaf node")

	midIndex := (tree.L + 1) / 2
	fmt.Printf("Calculated midIndex: %d\n", midIndex)

	newLeaf.keys = append(newLeaf.keys, node.keys[midIndex:]...)
	fmt.Printf("New leaf keys: %v\n", newLeaf.keys)

	node.keys = node.keys[:midIndex]
	fmt.Printf("Updated original node keys: %v\n", node.keys)

	newLeaf.next = node.next
	node.next = newLeaf
	fmt.Println("Updated next pointers")

	if node.parent == nil {
		fmt.Println("Node has no parent, creating new root")
		createNewRoot(tree, node, newLeaf)
	} else {
		fmt.Println("Inserting new leaf into parent")
		fmt.Printf("Parent keys before insertion: %v\n", node.parent.keys)
		insertIntoParent(tree, node.parent, newLeaf.keys[0], newLeaf)
		fmt.Printf("Parent keys after insertion: %v\n", node.parent.keys)
	}

	fmt.Println("Finished splitLeafNode function")
}

func createNewRoot(tree *BPTree, leftChild, rightChild *BPNode) {
	newRoot := &BPNode{
		keys:     []int{rightChild.keys[0]},
		children: []*BPNode{leftChild, rightChild},
		isLeaf:   false,
	}

	leftChild.parent = newRoot
	rightChild.parent = newRoot

	tree.root = newRoot
}

func insertIntoParent(tree *BPTree, parent *BPNode, key int, rightChild *BPNode) {
	insertKey(parent, key, rightChild)

	if len(parent.children) > tree.M {
		splitInternalNode(tree, parent)
	}
}

func insertKey(node *BPNode, key int, rightChild *BPNode) {
	i := 0
	for ; i < len(node.keys) && key > node.keys[i]; i++ {
	}

	node.keys = append(node.keys[:i], append([]int{key}, node.keys[i:]...)...)
	node.children = append(node.children[:i+1], append([]*BPNode{rightChild}, node.children[i+1:]...)...)
	rightChild.parent = node
}

func splitInternalNode(tree *BPTree, node *BPNode) {
	fmt.Println("Starting splitInternalNode function")
	newInternal := &BPNode{
		keys:   make([]int, 0),
		isLeaf: false,
	}

	midIndex := len(node.keys) / 2
	midKey := node.keys[midIndex]

	newInternal.keys = append(newInternal.keys, node.keys[midIndex+1:]...)
	newInternal.children = append(newInternal.children, node.children[midIndex+1:]...)

	for _, child := range newInternal.children {
		child.parent = newInternal
	}

	node.keys = node.keys[:midIndex] // Exclude the midKey
	node.children = node.children[:midIndex+1]

	if node.parent == nil {
		createNewRoot(tree, node, newInternal)
	} else {
		insertIntoParent(tree, node.parent, midKey, newInternal)
	}

	fmt.Println("Finished splitInternalNode function")
}

func (tree *BPTree) Delete(value int) {
	fmt.Printf("Deleting value %d from the B+ tree\n", value)

	node := findLeafNode(tree.root, value)
	fmt.Printf("Found leaf node with keys: %v\n", node.keys)

	deleteFromLeafNode(node, value)
	fmt.Printf("After deletion, leaf node keys: %v\n", node.keys)

	if len(node.keys) < (tree.L+1)/2 {
		handleLeafUnderflow(tree, node)
	}
}

func deleteFromLeafNode(node *BPNode, value int) {
	index := findKeyIndex(node.keys, value)
	if index != -1 {
		node.keys = append(node.keys[:index], node.keys[index+1:]...)
	}

	if node.parent != nil {
		updateParentKey(node.parent)
	}
}

func findKeyIndex(keys []int, value int) int {
	for i, key := range keys {
		if key == value {
			return i
		}
	}
	return -1
}

func handleLeafUnderflow(tree *BPTree, node *BPNode) {
	leftSibling := getLeftSibling(node)
	rightSibling := getRightSibling(node)

	if leftSibling != nil && len(leftSibling.keys) > (tree.L+1)/2 {
		borrowFromLeftSibling(node, leftSibling)
	} else if rightSibling != nil && len(rightSibling.keys) > (tree.L+1)/2 {
		borrowFromRightSibling(node, rightSibling)
	} else if leftSibling != nil {
		mergeWithSibling(tree, leftSibling, node)
	} else {
		mergeWithSibling(tree, node, rightSibling)
	}
}

func borrowFromLeftSibling(node, leftSibling *BPNode) {
	node.keys = append([]int{leftSibling.keys[len(leftSibling.keys)-1]}, node.keys...)
	leftSibling.keys = leftSibling.keys[:len(leftSibling.keys)-1]
	updateParentKey(node.parent)
}

func borrowFromRightSibling(node, rightSibling *BPNode) {
	node.keys = append(node.keys, rightSibling.keys[0])
	rightSibling.keys = rightSibling.keys[1:]
	updateParentKey(node.parent)
}

func mergeWithSibling(tree *BPTree, leftSibling, node *BPNode) {
	leftSibling.keys = append(leftSibling.keys, node.keys...)
	leftSibling.next = node.next

	if node.parent != nil {
		removeChildFromParent(tree, node.parent, node)
	}
}

func removeChildFromParent(tree *BPTree, parent, node *BPNode) {
	index := findChildIndex(parent.children, node)
	if index != -1 {
		parent.children = append(parent.children[:index], parent.children[index+1:]...)
		if index > 0 {
			parent.keys = append(parent.keys[:index-1], parent.keys[index:]...)
		} else {
			parent.keys = parent.keys[1:]
		}
	}

	if len(parent.children) < tree.MinChildrenInternalNodes() {
		handleInternalUnderflow(tree, parent)
	}
}
func findChildIndex(children []*BPNode, node *BPNode) int {
	for i, child := range children {
		if child == node {
			return i
		}
	}
	return -1
}

func handleInternalUnderflow(tree *BPTree, node *BPNode) {
	leftSibling := getLeftSibling(node)
	rightSibling := getRightSibling(node)

	if leftSibling != nil && len(leftSibling.children) > tree.MinChildrenInternalNodes() {
		borrowChildFromLeftSibling(node, leftSibling)
	} else if rightSibling != nil && len(rightSibling.children) > tree.MinChildrenInternalNodes() {
		borrowChildFromRightSibling(node, rightSibling)
	} else if leftSibling != nil {
		mergeInternalWithSibling(tree, leftSibling, node)
	} else {
		mergeInternalWithSibling(tree, node, rightSibling)
	}
}

func borrowChildFromLeftSibling(node, leftSibling *BPNode) {
	node.children = append([]*BPNode{leftSibling.children[len(leftSibling.children)-1]}, node.children...)
	leftSibling.children = leftSibling.children[:len(leftSibling.children)-1]
	updateParentKey(node.parent)
}
func borrowChildFromRightSibling(node, rightSibling *BPNode) {
	node.children = append(node.children, rightSibling.children[0])
	rightSibling.children = rightSibling.children[1:]
	updateParentKey(node.parent)
}

func mergeInternalWithSibling(tree *BPTree, leftSibling, node *BPNode) {
	leftSibling.keys = append(leftSibling.keys, node.keys...)
	leftSibling.children = append(leftSibling.children, node.children...)

	if node.parent != nil {
		removeChildFromParent(tree, node.parent, node)
	}
}

func getLeftSibling(node *BPNode) *BPNode {
	parent := node.parent
	if parent == nil {
		return nil
	}

	index := findChildIndex(parent.children, node)
	if index > 0 {
		return parent.children[index-1]
	}
	return nil
}

func getRightSibling(node *BPNode) *BPNode {
	parent := node.parent
	if parent == nil {
		return nil
	}

	index := findChildIndex(parent.children, node)
	if index < len(parent.children)-1 {
		return parent.children[index+1]
	}
	return nil
}

func updateParentKey(parent *BPNode) {
	for i := 1; i < len(parent.children); i++ {
		parent.keys[i-1] = findMinValue(parent.children[i])
	}
}
func findMinValue(node *BPNode) int {
	if node.isLeaf {
		return node.keys[0]
	}
	return findMinValue(node.children[0])
}

func (tree *BPTree) Search(value int) (*BPNode, int) {
	return searchInTree(tree.root, value)
}

func searchInTree(node *BPNode, value int) (*BPNode, int) {
	fmt.Printf("Searching for value %d in the B+ tree\n", value)

	if node.isLeaf {
		for i, key := range node.keys {
			if key == value {
				fmt.Printf("Value %d found at index %d in leaf node\n", value, i)
				return node, i
			}
		}
		fmt.Printf("Value %d not found in the B+ tree\n", value)
		return nil, -1
	}

	fmt.Printf("Searching in internal node with keys: %v\n", node.keys)
	for i, key := range node.keys {
		if value < key {
			fmt.Printf("Value %d is less than key %d, searching in child node\n", value, key)
			return searchInTree(node.children[i], value)
		}
	}
	fmt.Printf("Value %d is greater than all keys in the internal node, searching in the last child\n", value)
	return searchInTree(node.children[len(node.children)-1], value)
}
