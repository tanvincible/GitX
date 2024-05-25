package metadata

import (
	"GitX/models"
)

// Metadata represents the structure of the metadata file.
type Metadata struct {
	RepositoryName string          `json:"repository_name"`
	Description    string          `json:"description"`
	Branches       []string        `json:"branches"`
	Commits        []models.Commit `json:"commits"`
	// Add more fields as needed
}
