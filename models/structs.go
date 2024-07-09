package models

import "time"

// File represents a file in the repository.
type File struct {
	Path string
	Content string
}

// Branch represents a branch in the repository.
type Branch struct {
	Name   string
	Commit *Commit
}

// Repository represents the entire repository.
type Repository struct {
	Directory string
	Branches  []Branch
	HEAD      *Branch
}

// Reflog represents a reference log entry in the repository.
type Reflog struct {
	ID        string
	Author    string
	Timestamp time.Time
	Message   string
	// Add other necessary fields here
}

// GitXConfig represents your configuration settings.
type GitXConfig struct {
	UserName  string `toml:"user.name"`
	UserEmail string `toml:"user.email"`
	// Add other fields as needed
}

// IndexEntry represents an entry in the INDEX file.
type IndexEntry struct {
	Mode string `json:"mode"` // Mode field representing file mode
	Type string `json:"type"` // Type field representing object type (blob)
	Hash string `json:"hash"` // Hash field representing object hash
	Path string `json:"path"` // Path field representing file path
}

// IndexFile represents the INDEX file.
type IndexFile struct {
	Entries []*IndexEntry `json:"entries"`
}
