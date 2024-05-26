package vcs_operations

import (
	"GitX/models"
	"crypto/sha1"
	"encoding/hex"
)

type MerkleNode struct {
	Commit *models.Commit
	Left   *MerkleNode
	Right  *MerkleNode
}

func NewMerkleNode(left, right *MerkleNode, commit *models.Commit) *MerkleNode {
	node := &MerkleNode{}

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

func NewMerkleTree(commits []*models.Commit) *MerkleNode {
	var nodes []MerkleNode

	// Create leaf nodes
	for _, commit := range commits {
		nodes = append(nodes, *NewMerkleNode(nil, nil, commit))
	}

	for len(nodes) > 1 {
		var level []MerkleNode

		for i := 0; i < len(nodes); i += 2 {
			if i+1 < len(nodes) {
				level = append(level, *NewMerkleNode(&nodes[i], &nodes[i+1], nil))
			} else {
				level = append(level, nodes[i])
			}
		}

		nodes = level
	}

	return &nodes[0]
}
