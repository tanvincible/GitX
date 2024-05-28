package storage

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

// WriteFile writes data to a file.
func WriteFile(directory, filename string, data []byte) error {
	var filePath string
	if directory == "" {
		filePath = filename
	} else {
		// Create the directory if it doesn't exist
		if err := os.MkdirAll(directory, 0755); err != nil {
			return err
		}
		filePath = filepath.Join(directory, filename)
	}

	// Open the file for writing
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write data to the file
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// ReadFile reads data from a file.
func ReadFile(directory, filename string) ([]byte, error) {
	// Check if the file exists
	filePath := filepath.Join(directory, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, errors.New("file not found")
	}

	// Open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read data from the file
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// CreateStoragePath creates a storage path based on the base directory and hash value.
func CreateStoragePath(baseDir, hashValue string) (string, error) {
	dirName := hashValue[:2]
	fileName := hashValue[2:]
	dirPath := filepath.Join(baseDir, dirName)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return "", err
	}
	return filepath.Join(dirPath, fileName), nil
}

// StoreCompressedFile stores compressed data to a file.
func StoreCompressedFile(compressedData []byte, storagePath string) error {
	file, err := os.Create(storagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(compressedData)
	if err != nil {
		return err
	}

	return nil
}
