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
	var metadata metadata.Metadata

	// Read the JSON data from the metadata file
	filePath := filepath.Join(directory, "metadata.json")
	file, err := os.Open(filePath)
	if err != nil {
		return metadata, err
	}
	defer file.Close()

	// Decode the JSON data into metadata structure
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&metadata); err != nil {
		return metadata, err
	}

	return metadata, nil
}

// UpdateMetadata updates the metadata file with new data.
func UpdateMetadata(metadataFile, filePath, hashValue string, newCommit models.Commit) error {
	// Read existing metadata
	metadata, err := ReadMetadata(filepath.Dir(metadataFile))
	if err != nil {
		return err
	}

	// Update metadata with new commit information
	commit := models.Commit{
		ID:      hashValue,
		Author:  newCommit.Author,
		Message: newCommit.Message,
		// Additional fields as needed
	}
	metadata.Commits = append(metadata.Commits, commit)

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
