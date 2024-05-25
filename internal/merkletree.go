package main

import (
	"GitX/models"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

type Node struct {
	Commit *models.Commit
	Left   *Node
	Right  *Node
}

func NewNode(left, right *Node, commit *models.Commit) *Node {
	node := &Node{}

	if left == nil && right == nil {
		hash := sha1.Sum([]byte(commit.ID))
		node.Commit = commit
		node.Commit.ID = hex.EncodeToString(hash[:])
	} else {
		prevHashes := append([]byte(left.Commit.ID), []byte(right.Commit.ID)...)
		hash := sha1.Sum(prevHashes)
		node.Commit = &models.Commit{ID: hex.EncodeToString(hash[:])}
	}

	node.Left = left
	node.Right = right

	return node
}

func NewMerkleTree(commits []*models.Commit) *Node {
	var nodes []Node

	// Create leaf nodes
	for _, commit := range commits {
		nodes = append(nodes, *NewNode(nil, nil, commit))
	}

	for len(nodes) > 1 {
		var level []Node

		for i := 0; i < len(nodes); i += 2 {
			if i+1 < len(nodes) {
				level = append(level, *NewNode(&nodes[i], &nodes[i+1], nil))
			} else {
				level = append(level, nodes[i])
			}
		}

		nodes = level
	}

	return &nodes[0]
}

func main() {
	// Sample commits
	commits := []*models.Commit{
		{ID: "Commit 1"},
		{ID: "Commit 2"},
		{ID: "Commit 3"},
		{ID: "Commit 4"},
	}

	// Create a new Merkle tree from the commits
	root := NewMerkleTree(commits)

	// Print the Merkle tree
	printMerkleTree(root, 0)
}

func printMerkleTree(node *Node, level int) {
	if node == nil {
		return
	}

	format := ""
	for i := 0; i < level; i++ {
		format += "\t"
	}

	format += "---[ "
	level++
	fmt.Printf(format+"%s\n", node.Commit.ID)

	printMerkleTree(node.Left, level)
	printMerkleTree(node.Right, level)
}
