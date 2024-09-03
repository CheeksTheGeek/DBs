#include <stdio.h>
#include <stdlib.h>
#include "b+tree.h"

/**
 * @brief 
 * 
 * @param degree 
 * @return BPlusTree* 
 */
BPlusTree *createBPlusTree(int degree) {
    BPlusTree *tree = (BPlusTree *)malloc(sizeof(BPlusTree));
    tree->root = NULL;
    tree->degree = degree;
    return tree;
}

/**
 * @brief Create a Node object
 * 
 * @param degree The degree of the B+ tree
 * @param isLeaf 1 if leaf node, 0 otherwise
 * @return BPlusTreeNode* 
 */
BPlusTreeNode *createNode(int degree, int isLeaf) {
    BPlusTreeNode *newNode = (BPlusTreeNode *)malloc(sizeof(BPlusTreeNode));
    newNode->isLeaf = isLeaf;
    newNode->numKeys = 0;
    for (int i = 0; i < degree + 1; i++) {  // Fixed the bounds to degree + 1
        newNode->children[i] = NULL;
    }
    newNode->next = NULL;
    return newNode;
}

/**
 * @brief Insert a key into the B+ tree
 * 
 * @param tree 
 * @param key 
 */
void insertKey(BPlusTree *tree, int key) {
    printf("Inserting key: %d\n", key);
    if (!tree) {
        printf("Tree is NULL! Aborting insert.\n");
        return;
    }

    if (tree->root == NULL) {
        printf("Root is NULL, creating a new root...\n");
        tree->root = createNode(tree->degree, 1);
        tree->root->keys[0] = key;
        tree->root->numKeys = 1;
        printf("Root created with key: %d\n", key);
    } else {
        BPlusTreeNode *root = tree->root;
        if (root->numKeys == 2 * tree->degree - 1) {
            printf("Root is full, creating new root...\n");
            BPlusTreeNode *newRoot = createNode(tree->degree, 0);
            newRoot->children[0] = root;
            splitNode(tree, newRoot, 0);
            tree->root = newRoot;
            insertNonFull(tree, newRoot, key);
        } else {
            insertNonFull(tree, root, key);
        }
    }
}

void splitNode(BPlusTree *tree, BPlusTreeNode *parent, int index) {
    printf("Splitting node at index %d...\n", index);

    BPlusTreeNode *node = parent->children[index];
    if (node == NULL) {
        printf("Error: Node to split is NULL.\n");
        return;
    }

    int degree = tree->degree;
    int mid = (degree - 1) / 2;

    BPlusTreeNode *newNode = createNode(degree, node->isLeaf);
    if (newNode == NULL) {
        printf("Error: Could not create new node during split.\n");
        return;
    }

    // Move second half of keys from the old node to the new node
    for (int i = 0; i < degree - 1 - mid; i++) {
        newNode->keys[i] = node->keys[i + mid + 1];
    }
    
    if (!node->isLeaf) {
        // Move corresponding children for internal nodes
        for (int i = 0; i < degree - mid; i++) {
            newNode->children[i] = node->children[i + mid + 1];
            node->children[i + mid + 1] = NULL;
        }
    }
    
    newNode->numKeys = degree - 1 - mid;
    node->numKeys = mid;

    // Insert new key into parent
    for (int i = parent->numKeys; i > index; i--) {
        parent->children[i + 1] = parent->children[i];
        parent->keys[i] = parent->keys[i - 1];
    }
    parent->keys[index] = node->keys[mid];
    parent->children[index + 1] = newNode;
    parent->numKeys++;

    // For leaf nodes, maintain the linked list
    if (node->isLeaf) {
        newNode->next = node->next;
        node->next = newNode;
    }

    printf("Split complete: Parent keys: ");
    for (int i = 0; i < parent->numKeys; i++) {
        printf("%d ", parent->keys[i]);
    }
    printf("\n");
}

void insertNonFull(BPlusTree *tree, BPlusTreeNode *node, int key) {
    printf("Inserting %d into non-full node...\n", key);

    int i = node->numKeys - 1;
    if (node->isLeaf) {
        while (i >= 0 && key < node->keys[i]) {
            node->keys[i + 1] = node->keys[i];
            i--;
        }
        node->keys[i + 1] = key;
        node->numKeys++;
        printf("Inserted key %d into leaf node.\n", key);
    } else {
        while (i >= 0 && key < node->keys[i]) {
            i--;
        }
        i++;

        if (node->children[i]->numKeys == 2 * tree->degree - 1) {
            splitNode(tree, node, i);
            if (key > node->keys[i]) {
                i++;
            }
        }
        insertNonFull(tree, node->children[i], key);
    }
}



BPlusTreeNode *search(BPlusTree *tree, int key) {
    return searchNode(tree->root, key);
}

BPlusTreeNode *searchNode(BPlusTreeNode *node, int key) {
    int i = 0;
    while (i < node->numKeys && key > node->keys[i]) {
        i++;
    }

    if (i < node->numKeys && key == node->keys[i]) {
        return node; // Found key
    }

    if (node->isLeaf) {
        return NULL; // Key not found
    } else {
        return searchNode(node->children[i], key);
    }
}

void printTree(BPlusTree *tree) {
    if (tree->root != NULL) {
        printNode(tree->root, 0);
    } else {
        printf("Empty Tree\n");
    }
}

void printNode(BPlusTreeNode *node, int level) {
    printf("Level %d: ", level);
    for (int i = 0; i < node->numKeys; i++) {
        printf("%d ", node->keys[i]);
    }
    printf("\n");

    if (!node->isLeaf) {
        for (int i = 0; i <= node->numKeys; i++) {
            printNode(node->children[i], level + 1);
        }
    }
}
