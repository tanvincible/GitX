package metadata_operations

import (
	"GitX/internal/metadata"
	"GitX/models"
	"encoding/json"
	"os"
	"path/filepath"
)

// WriteMetadata writes the metadata to a file.
func WriteMetadata(metadata metadata.Metadata, directory string) error {

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(directory, 0755); err != nil {
		return err
	}

	// Marshal the metadata to JSON
	data, err := json.MarshalIndent(metadata, "", "    ")
	if err != nil {
		return err
	}

	// Write the JSON data to the metadata file
	filePath := filepath.Join(directory, "metadata.json")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// ReadMetadata reads the metadata from a file.
func ReadMetadata(directory string) (metadata.Metadata, error) {
	var meta metadata.Metadata // Renamed variable to avoid conflict

	// Read the JSON data from the metadata file
	filePath := filepath.Join(directory, "metadata.json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// If the metadata file doesn't exist, return a new Metadata instance
			return metadata.Metadata{
				RepositoryName: "",
				Description:    "",
				Branches:       []string{},
				Commits:        []models.Commit{},
			}, nil
		}
		return meta, err // Updated variable name
	}

	// Check if the file is empty
	if len(data) == 0 {
		return metadata.Metadata{
			RepositoryName: "",
			Description:    "",
			Branches:       []string{},
			Commits:        []models.Commit{},
		}, nil
	}

	// Decode the JSON data into the Metadata structure
	err = json.Unmarshal(data, &meta) // Updated variable name
	if err != nil {
		return meta, err // Updated variable name
	}

	return meta, nil // Updated variable name
}
// UpdateMetadata updates the metadata file with new data.
func UpdateMetadata(metadataFile string, newCommit models.Commit, filePath string, hashValue string) error {
	// Read existing metadata
	metadata, err := ReadMetadata(filepath.Dir(metadataFile))
	if err != nil {
		return err
	}

	// If filePath and hashValue are provided, update the commit with them
	if filePath != "" && hashValue != "" {
		newCommit.Files = map[string]string{filePath: hashValue}
	}

	// Update metadata with new commit information
	metadata.Commits = append(metadata.Commits, newCommit)

	// Write updated metadata back to file
	if err := WriteMetadata(metadata, filepath.Dir(metadataFile)); err != nil {
		return err
	}

	return nil
}

// GetTrackedFiles retrieves the tracked files from the metadata.
func GetTrackedFiles(metadataFile string) (map[string]string, error) {
	metadata, err := ReadMetadata(filepath.Dir(metadataFile))
	if err != nil {
		return nil, err
	}

	trackedFiles := make(map[string]string)
	for _, commit := range metadata.Commits {
		for path, content := range commit.Files {
			trackedFiles[path] = content
		}
	}

	return trackedFiles, nil
}
