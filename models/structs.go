package models

// File represents a file in the repository.
type File struct {
	Path string
	Hash string
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
