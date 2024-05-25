package models

import (
	"crypto/sha1"
	"fmt"
	"os"
)

// Blob represents a file in the repository with its content.
type Blob struct {
	ID      string // Unique identifier, typically a hash of the content
	Content string // The content of the file
}

// NewBlob creates a new Blob from a file.
func NewBlob(filePath string) (*Blob, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	blob := &Blob{
		ID:      calculateHash(content),
		Content: string(content),
	}

	return blob, nil
}

// calculateHash calculates the SHA-1 hash of the content.
func calculateHash(content []byte) string {
	hasher := sha1.New()
	hasher.Write(content)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
