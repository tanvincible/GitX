package file_operations

import (
	"GitX/internal/compression"
	"GitX/internal/hash"
	"GitX/internal/storage"
	"GitX/models"
	"GitX/utils/metadata_operations"
	"GitX/utils/vcs_operations"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func InitHandler(directory string) {
	// Create repository directory
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		log.Fatalf("Error creating repository directory: %v", err)
	}

	// Create .gitx directory inside the repository directory
	gitxDir := filepath.Join(directory, ".gitx")
	if err := os.MkdirAll(gitxDir, os.ModePerm); err != nil {
		log.Fatalf("Error creating .gitx directory: %v", err)
	}

	// Create metadata file
	metadataFile := filepath.Join(gitxDir, "metadata.json")
	if _, err := os.Create(metadataFile); err != nil {
		log.Fatalf("Error creating metadata file: %v", err)
	}

	// Create HEAD file
	headFile := filepath.Join(gitxDir, "HEAD")
	if err := os.WriteFile(headFile, []byte("refs/heads/main"), 0644); err != nil {
		log.Fatalf("Error creating HEAD file: %v", err)
	}

	// Create refs/heads directory
	refsHeadsDir := filepath.Join(gitxDir, "refs", "heads")
	if err := os.MkdirAll(refsHeadsDir, os.ModePerm); err != nil {
		log.Fatalf("Error creating refs/heads directory: %v", err)
	}

	// Create main branch file
	mainBranchFile := filepath.Join(refsHeadsDir, "main")
	if _, err := os.Create(mainBranchFile); err != nil {
		log.Fatalf("Error creating main branch file: %v", err)
	}

	// Create commits directory
	commitsDir := filepath.Join(gitxDir, "commits")
	if err := os.MkdirAll(commitsDir, os.ModePerm); err != nil {
		log.Fatalf("Error creating commits directory: %v", err)
	}

	// Set up ignore file
	ignoreFile := filepath.Join(directory, ".gitxignore")
	if _, err := os.Create(ignoreFile); err != nil {
		log.Fatalf("Error creating ignore file: %v", err)
	}

	fmt.Printf("Initialized empty repository in %s\n", directory)
}

func AddHandler(filePath string, stagingArea map[string]string) {
	// Step 1: Calculate SHA-1 Hash
	hashValue, err := hash.SHA1Hash(filePath)
	if err != nil {
		log.Fatalf("Error calculating hash for file %s: %v", filePath, err)
	}

	// Step 2: Add file to staging area
	stagingArea[filePath] = hashValue

	fmt.Printf("File %s added to staging area\n", filePath)
}

func CommitHandler(message string, stagingArea map[string]string) {
	// Step 1: For each file in the staging area
	for filePath, hashValue := range stagingArea {
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

		// Step 5: Create Commit Object
		newCommit := models.Commit{
			ID:        hashValue,
			Parent:    nil, // Set parent commit if applicable
			Tree:      nil, // Set tree reference
			Message:   message,
			Author:    "",         // Set author information
			Timestamp: time.Now(), // Set commit timestamp
			// Additional fields if needed
		}
		// Step 7: Update Metadata
		metadataFile := "metadata.json"
		if err := metadata_operations.UpdateMetadata(metadataFile, filePath, hashValue, newCommit); err != nil {
			log.Fatalf("Error updating metadata: %v", err)
		}
		// Step 8: Update HEAD to point to the new commit
		vcs_operations.UpdateHEAD(hashValue)
	}

	// Step 9: Clear the staging area
	for k := range stagingArea {
		delete(stagingArea, k)
	}

	fmt.Printf("Commit created with message: %s\n", message)
}

func StatusHandler(stagingArea map[string]string) {
	// Step 1: Retrieve tracked files from metadata
	trackedFiles, err := metadata_operations.GetTrackedFiles("metadata.json")
	if err != nil {
		log.Fatalf("Error retrieving tracked files: %v", err)
	}

	// Step 2: Get list of files in the working directory
	workingDirFiles, err := getAllFilesInDir(".")
	if err != nil {
		log.Fatalf("Error retrieving files from working directory: %v", err)
	}

	// Step 3: Compare files
	fmt.Println("Changes to be committed:")
	for filePath, hashValue := range stagingArea {
		if _, ok := trackedFiles[filePath]; !ok {
			fmt.Printf("\tnew file: %s\n", filePath)
		} else {
			if hashValue != trackedFiles[filePath] {
				fmt.Printf("\tmodified: %s\n", filePath)
			}
		}
	}

	fmt.Println("Changes not staged for commit:")
	for _, file := range workingDirFiles {
		if _, ok := stagingArea[file]; !ok {
			if hashValue, err := hash.SHA1Hash(file); err == nil {
				if trackedHash, ok := trackedFiles[file]; ok {
					if hashValue != trackedHash {
						fmt.Printf("\tmodified: %s\n", file)
					}
				} else {
					fmt.Printf("\tuntracked: %s\n", file)
				}
			} else {
				log.Printf("Error hashing file %s: %v", file, err)
			}
		}
	}
}

func getAllFilesInDir(dirPath string) ([]string, error) {
	var files []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
