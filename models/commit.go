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
