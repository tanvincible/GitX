package models

import (
	"time"
)

type Commit struct {
	ID        string
	Parent    []*Commit
	Tree      *Tree
	Message   string
	Author    string
	Timestamp time.Time
	Files     map[string]string
	// Additional fields
	Committer    string
	GPGSignature string
}

// Additional methods

// GetShortID returns the first 7 characters of the commit ID.
/*
func (c *Commit) GetShortID() string {
	return c.ID[:7]
}

// VerifySignature verifies the GPG signature of the commit.
func (c *Commit) VerifySignature() bool {
	// Implementation of signature verification
}

// GetChangedFiles returns a list of files changed in this commit.
func (c *Commit) GetChangedFiles() []string {
	// Implementation of getting changed files
}
*/
