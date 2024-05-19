package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"GitX/internal/compression"
	"GitX/internal/hash"

	// "GitX/internal/metadata"
	"GitX/internal/storage"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Usage: gitx <command> [init, commit]")
		os.Exit(1)
	}

	command := args[0]
	switch command {
	case "init":
		if len(args) != 2 {
			fmt.Println("Usage: gitx init <directory>")
			os.Exit(1)
		}
		InitHandler(args[1])
	case "commit":
		if len(args) != 3 {
			fmt.Println("Usage: gitx commit <file-path> <message>")
			os.Exit(1)
		}
		CommitHandler(args[1], args[2])
	default:
		fmt.Printf("gitx: %s is not a valid command\n", command)
		os.Exit(1)
	}
}

func InitHandler(directory string) {
	// Create repository directory
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		log.Fatalf("Error creating repository directory: %v", err)
	}

	// Create metadata file
	metadataFile := filepath.Join(directory, "metadata.json")
	if _, err := os.Create(metadataFile); err != nil {
		log.Fatalf("Error creating metadata file: %v", err)
	}

	// Create default branch
	//

	// Set up ignore file
	ignoreFile := filepath.Join(directory, ".gitignore")
	if _, err := os.Create(ignoreFile); err != nil {
		log.Fatalf("Error creating ignore file: %v", err)
	}

	// Initialize empty commit object
	//

	fmt.Printf("Initialized repository in %s\n", directory)
}

func CommitHandler(filePath, message string) {
	// Step 1: Calculate SHA-1 Hash
	hashValue, err := hash.SHA1Hash(filePath)
	if err != nil {
		log.Fatalf("Error calculating hash for file %s: %v", filePath, err)
	}

	// Step 2: Create Storage Path
	baseDir := "storage"
	storagePath, err := storage.CreateStoragePath(baseDir, hashValue)
	if err != nil {
		log.Fatalf("Error creating storage path: %v", err)
	}

	// Step 3: Compress File Content
	compressedData, err := compression.CompressFile(filePath)
	if err != nil {
		log.Fatalf("Error compressing file: %v", err)
	}

	// Step 4: Store Compressed File
	if err := storage.StoreCompressedFile(compressedData, storagePath); err != nil {
		log.Fatalf("Error storing compressed file: %v", err)
	}

	// Step 5: Update Metadata
	/*
		metadataFile := "metadata.json"
		if err := metadata.UpdateMetadata(metadataFile, filePath, hashValue); err != nil {
			log.Fatalf("Error updating metadata: %v", err)
		}
	*/

	// Step 6: Create Commit Object
	// (You can implement this step based on your VCS requirements)
	fmt.Printf("Commit created for file %s with message: %s\n", filePath, message)
}
