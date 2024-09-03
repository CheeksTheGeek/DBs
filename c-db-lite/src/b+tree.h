#pragma once
#define DEGREE 4
#define MAX_KEYS (2 * DEGREE - 1)

typedef struct BPlusTreeNode {
    int keys[MAX_KEYS];                     // Array of keys
    struct BPlusTreeNode *children[MAX_KEYS + 1]; // Children pointers
    int isLeaf;                             // 1 if leaf, 0 otherwise
    int numKeys;                            // Current number of keys
    struct BPlusTreeNode *next;             // Pointer to next leaf node (if leaf)
} BPlusTreeNode;

typedef struct BPlusTree {
    BPlusTreeNode *root;                    // Root node
    int degree;                             // Degree of the tree
} BPlusTree;

// Function declarations
BPlusTree *createBPlusTree(int degree);
BPlusTreeNode *createNode(int degree, int isLeaf);
void splitNode(BPlusTree *tree, BPlusTreeNode *parent, int index);
void insertKey(BPlusTree *tree, int key);
void insertNonFull(BPlusTree *tree, BPlusTreeNode *node, int key);
void deleteKey(BPlusTree *tree, int key);
BPlusTreeNode *search(BPlusTree *tree, int key);
BPlusTreeNode *searchNode(BPlusTreeNode *node, int key);
void printTree(BPlusTree *tree);
void printNode(BPlusTreeNode *node, int level);
