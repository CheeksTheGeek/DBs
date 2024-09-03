#include <stdio.h>
#include <stdlib.h>
#include "b+tree.h"

/**
 * @brief Test function for createBPlusTree
 * 
 */
void testCreateBPlusTree() {
    printf("Testing createBPlusTree...\n");
    BPlusTree *tree = createBPlusTree(4);
    if (tree && tree->degree == 4 && tree->root == NULL) {
        printf("createBPlusTree passed!\n");
    } else {
        printf("createBPlusTree failed!\n");
    }
}

/**
 * @brief Test function for createNode
 * 
 */
void testCreateNode() {
    printf("Testing createNode...\n");
    BPlusTreeNode *node = createNode(4, 1);
    if (node && node->isLeaf == 1 && node->numKeys == 0) {
        printf("createNode passed!\n");
    } else {
        printf("createNode failed!\n");
    }
}


/**
 * @brief Test function for insertKey
 * 
 * @param tree 
 */
void testInsertKey(BPlusTree *tree) {
    printf("Testing insertKey...\n");
    insertKey(tree, 10);
    insertKey(tree, 20);
    insertKey(tree, 5);
    insertKey(tree, 6);
    insertKey(tree, 12);
    insertKey(tree, 30);
    insertKey(tree, 7);
    insertKey(tree, 17);

    printTree(tree);  // Display the tree

    BPlusTreeNode *result = search(tree, 12);
    if (result != NULL) {
        printf("Key 12 found during testInsertKey.\n");
    } else {
        printf("Key 12 not found during testInsertKey.\n");
    }
}
/**
 * @brief Main function to run B+ Tree tests
 * 
 * @return int 
 */
int main() {
    printf("Starting B+ Tree tests...\n");
    // Test createBPlusTree
    testCreateBPlusTree();
    // Create a tree for further tests
    BPlusTree *tree = createBPlusTree(4);
    // Test createNode
    testCreateNode();
    // Test insertKey with extensive logging
    testInsertKey(tree);
    printf("B+ Tree tests completed.\n");
    return 0;
}