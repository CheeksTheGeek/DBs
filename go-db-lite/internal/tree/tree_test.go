package tree

import (
	"testing"
)

// Test insertion of elements into the B+ tree
func TestInsert(t *testing.T) {
	tree := CreateTree(4, 4)

	// Insert a sequence of elements
	values := []int{20, 18, 22, 28, 25, 41, 23, 42, 53, 35, 62, 63, 84, 73, 99}

	for _, value := range values {
		tree.Insert(value)
	}

	// Check the structure and values in the tree
	root := tree.root
	if len(root.keys) != 1 || root.keys[0] != 25 {
		t.Errorf("Root node key is incorrect: got %v, want [25]", root.keys)
	}

	if len(root.children) != 2 {
		t.Errorf("Root node should have 2 children: got %d", len(root.children))
	}

	leftChild := root.children[0]
	rightChild := root.children[1]

	if len(leftChild.keys) != 2 || leftChild.keys[0] != 23 || leftChild.keys[1] != 41 {
		t.Errorf("Left child node keys are incorrect: got %v, want [23, 41]", leftChild.keys)
	}

	if len(rightChild.keys) != 2 || rightChild.keys[0] != 63 || rightChild.keys[1] != 84 {
		t.Errorf("Right child node keys are incorrect: got %v, want [63, 84]", rightChild.keys)
	}
}

// Test deletion of elements from the B+ tree
func TestDelete(t *testing.T) {
	tree := CreateTree(3, 4)

	// Insert a sequence of elements
	values := []int{20, 18, 22, 28, 25, 41, 23, 42, 53, 35, 62, 63, 84, 73, 99}
	for _, value := range values {
		tree.Insert(value)
	}

	// Delete some elements
	toDelete := []int{41, 28, 62, 99, 20, 63}
	for _, value := range toDelete {
		tree.Delete(value)
	}

	// Check the structure and values in the tree after deletions
	root := tree.root
	if len(root.keys) != 1 || root.keys[0] != 25 {
		t.Errorf("Root node key is incorrect after deletion: got %v, want [25]", root.keys)
	}

	leftChild := root.children[0]
	rightChild := root.children[1]

	if len(leftChild.keys) != 2 || leftChild.keys[0] != 23 || leftChild.keys[1] != 35 {
		t.Errorf("Left child node keys are incorrect after deletion: got %v, want [23, 35]", leftChild.keys)
	}

	if len(rightChild.keys) != 2 || rightChild.keys[0] != 53 || rightChild.keys[1] != 84 {
		t.Errorf("Right child node keys are incorrect after deletion: got %v, want [53, 84]", rightChild.keys)
	}
}

// Test tree rebalancing during multiple insertions and deletions
func TestRebalancing(t *testing.T) {
	tree := CreateTree(3, 4)

	// Insert elements to cause multiple splits
	values := []int{15, 25, 35, 45, 55, 65, 75, 85, 95, 5, 35, 65, 75, 85, 95}
	for _, value := range values {
		tree.Insert(value)
	}

	// Delete elements to cause multiple merges
	toDelete := []int{15, 25, 35, 45, 55, 65, 75, 85, 95}
	for _, value := range toDelete {
		tree.Delete(value)
	}

	// Check tree root
	root := tree.root
	if len(root.keys) != 1 || root.keys[0] != 65 {
		t.Errorf("Root node key is incorrect after rebalancing: got %v, want [65]", root.keys)
	}

	// Check children
	if len(root.children) != 2 {
		t.Errorf("Root should have 2 children after rebalancing: got %d", len(root.children))
	}

	leftChild := root.children[0]
	rightChild := root.children[1]

	if len(leftChild.keys) != 2 || leftChild.keys[0] != 5 || leftChild.keys[1] != 35 {
		t.Errorf("Left child node keys are incorrect after rebalancing: got %v, want [5, 35]", leftChild.keys)
	}

	if len(rightChild.keys) != 2 || rightChild.keys[0] != 75 || rightChild.keys[1] != 95 {
		t.Errorf("Right child node keys are incorrect after rebalancing: got %v, want [75, 95]", rightChild.keys)
	}
}

// Test merging of nodes during deletion
func TestMerging(t *testing.T) {
	tree := CreateTree(3, 4)

	// Insert elements
	values := []int{10, 20, 30, 40, 50, 60, 70, 80, 90}
	for _, value := range values {
		tree.Insert(value)
	}

	// Delete elements to trigger merges
	toDelete := []int{80, 70, 60, 50, 40}
	for _, value := range toDelete {
		tree.Delete(value)
	}

	// Check the structure of the tree after merges
	root := tree.root
	if len(root.keys) != 1 || root.keys[0] != 20 {
		t.Errorf("Root node key is incorrect after merging: got %v, want [20]", root.keys)
	}

	if len(root.children) != 2 {
		t.Errorf("Root should have 2 children after merging: got %d", len(root.children))
	}

	leftChild := root.children[0]
	rightChild := root.children[1]

	if len(leftChild.keys) != 1 || leftChild.keys[0] != 10 {
		t.Errorf("Left child node keys are incorrect after merging: got %v, want [10]", leftChild.keys)
	}

	if len(rightChild.keys) != 2 || rightChild.keys[0] != 30 || rightChild.keys[1] != 90 {
		t.Errorf("Right child node keys are incorrect after merging: got %v, want [30, 90]", rightChild.keys)
	}
}

// Test empty tree case
func TestEmptyTree(t *testing.T) {
	tree := CreateTree(3, 4)

	if tree.root == nil {
		t.Errorf("Root should not be nil for an empty tree")
	}

	if len(tree.root.keys) != 0 {
		t.Errorf("Root keys should be empty in an empty tree: got %v", tree.root.keys)
	}
}

// Test edge cases
func TestEdgeCases(t *testing.T) {
	tree := CreateTree(3, 4)

	// Insert the same element multiple times
	for i := 0; i < 5; i++ {
		tree.Insert(10)
	}

	if len(tree.root.keys) != 1 || tree.root.keys[0] != 10 {
		t.Errorf("Root node should contain only one key 10 after multiple insertions of the same element: got %v", tree.root.keys)
	}

	if len(tree.root.children) > 0 {
		t.Errorf("Root node should not have children when inserting the same element multiple times")
	}

	// Delete the element
	tree.Delete(10)
	if len(tree.root.keys) != 0 {
		t.Errorf("Root node should be empty after deleting all instances of the element: got %v", tree.root.keys)
	}
}

// Test borrowing from sibling
func TestBorrowing(t *testing.T) {
	tree := CreateTree(3, 4)

	// Insert elements to create a scenario for borrowing
	values := []int{10, 20, 30, 40, 50, 60, 70, 80}
	for _, value := range values {
		tree.Insert(value)
	}

	// Delete to trigger borrowing
	tree.Delete(80)
	tree.Delete(70)

	root := tree.root

	if len(root.keys) != 1 || root.keys[0] != 40 {
		t.Errorf("Root node key is incorrect after borrowing: got %v, want [40]", root.keys)
	}

	leftChild := root.children[0]
	rightChild := root.children[1]

	if len(leftChild.keys) != 2 || leftChild.keys[0] != 10 || leftChild.keys[1] != 30 {
		t.Errorf("Left child node keys are incorrect after borrowing: got %v, want [10, 30]", leftChild.keys)
	}

	if len(rightChild.keys) != 2 || rightChild.keys[0] != 50 || rightChild.keys[1] != 60 {
		t.Errorf("Right child node keys are incorrect after borrowing: got %v, want [50, 60]", rightChild.keys)
	}
}
