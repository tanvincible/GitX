package models

// TreeEntry represents an entry in a tree, which can be either another tree (directory) or a blob (file).
type TreeEntry struct {
	Name string // The name of the entry
	Mode string // The file permissions mode
	ID   string // The SHA-1 hash of the entry
	Type string // The type of the entry: "blob" or "tree"
}

// Tree represents a directory in the repository.
type Tree struct {
	ID      string      // The SHA-1 hash of the tree
	Entries []TreeEntry // The entries in the directory
}
