package file_operations

import (
	"GitX/internal/hash"
	"GitX/models"
	"GitX/utils/metadata_operations"
	"GitX/utils/vcs_operations"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// InitHandler initializes a new GitX repository by creating the necessary directories and files.
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

	// Create objects directory
	objectsDir := filepath.Join(gitxDir, "objects")
	if err := os.MkdirAll(objectsDir, os.ModePerm); err != nil {
		log.Fatalf("Error creating objects directory: %v", err)
	}

	// Set up ignore file
	ignoreFile := filepath.Join(directory, ".gitxignore")
	if _, err := os.Create(ignoreFile); err != nil {
		log.Fatalf("Error creating ignore file: %v", err)
	}

	// Create config file
	configFile := filepath.Join(gitxDir, "config")
	if _, err := os.Create(configFile); err != nil {
		log.Fatalf("Error creating config file: %v", err)
	}

	// Create description file
	descriptionFile := filepath.Join(gitxDir, "description")
	descriptionContent := []byte("Unnamed repository; edit this file to name the repository.\n")
	if err := os.WriteFile(descriptionFile, descriptionContent, 0644); err != nil {
		log.Fatalf("Error creating description file: %v", err)
	}

	fmt.Printf("Initialized empty repository in %s\n", directory)
}

// AddHandler adds a file to the staging area, handling paths in a cross-platform manner.
func AddHandler(repoRoot, filePath string, stagingArea map[string]string) error {
	// Ensure the repoRoot is an absolute path
	absRepoRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return fmt.Errorf("unable to resolve repository root: %w", err)
	}

	// Check if filePath is empty
	if filePath == "" {
		return fmt.Errorf("file path is empty")
	}

	// Join the repoRoot with the filePath to get the absolute path to the file
	absFilePath := filepath.Join(absRepoRoot, filePath)

	// Check if the file exists
	if _, err := os.Stat(absFilePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", absFilePath)
	}

	// Calculate the SHA-1 hash of the file and store the blob in the object database
	hashValue, err := hash.SHA1Hash(absFilePath)
	if err != nil {
		return fmt.Errorf("error calculating hash for file %s: %w", absFilePath, err)
	}

	// Normalize the file path to use forward slashes
	normalizedPath := filepath.ToSlash(filePath)

	// Add the file and its hash to the staging area
	stagingArea[normalizedPath] = hashValue

	return nil
}

// CommitHandler creates a commit object, compresses the file content, stores the compressed file, updates metadata, and updates the HEAD reference.
func CommitHandler(message string, stagingArea *map[string]string) {

	// Create the commits directory if it doesn't exist
	commitsDir := ".gitx/commits"
	if _, err := os.Stat(commitsDir); os.IsNotExist(err) {
		os.MkdirAll(commitsDir, os.ModePerm)
	}

	// Retrieve the current HEAD commit to set as the parent for the new commit
	parentCommitHash := vcs_operations.GetCurrentHeadCommit()
	var parentCommit *models.Commit
	var err error

	// Check if there is an existing parent commit
	if parentCommitHash != "" {
		parentCommit, err = vcs_operations.GetCommitByHash(parentCommitHash)
		if err != nil {
			log.Fatalf("Error retrieving parent commit: %v", err)
		}
	} else {
		// If there is no parent commit, create an initial commit for the main branch
		initialCommit := createInitialCommit()

		// Write the initial commit to the commits directory
		initialCommitData, err := json.Marshal(initialCommit)
		if err != nil {
			log.Fatalf("Error serializing initial commit data: %v", err)
		}
		initialCommitFilePath := filepath.Join(commitsDir, initialCommit.ID)
		if err := os.WriteFile(initialCommitFilePath, initialCommitData, 0644); err != nil {
			log.Fatalf("Error writing initial commit file: %v", err)
		}

		// Update HEAD to point to the initial commit
		if err := vcs_operations.UpdateHEAD(initialCommit.ID); err != nil {
			log.Fatalf("Error updating HEAD with initial commit: %v", err)
		}

		// Set the initial commit as the parent
		parentCommit = &initialCommit
	}

	// Create a tree object from the staging area files
	tree, err := vcs_operations.CreateTreeFromStagingArea(*stagingArea)
	if err != nil {
		log.Fatalf("Error creating tree from staging area: %v", err)
	}

	// Create a new commit object
	newCommit := models.Commit{
		ID:        "", // This will be generated based on the tree hash and parent
		Parent:    []*models.Commit{},
		Tree:      tree,
		Message:   message,
		Author:    vcs_operations.GetCurrentUser(), // Implement a function to get the current user
		Timestamp: time.Now(),
	}

	// If there is a parent commit, set it
	if parentCommit != nil {
		newCommit.Parent = append(newCommit.Parent, parentCommit)
	}

	// Generate the commit ID based on the tree hash, parent, author, committer, and timestamp
	newCommit.ID, err = vcs_operations.GenerateCommitID(newCommit.Tree, newCommit.Parent, newCommit.Message, newCommit.Author, newCommit.Timestamp)
	if err != nil {
		log.Fatalf("Error generating commit ID: %v", err)
	}

	// Serialize the commit object to JSON
	commitData, err := json.Marshal(newCommit)
	if err != nil {
		log.Fatalf("Error serializing commit data: %v", err)
	}

	// Write Commit Object to File
	commitFilePath := filepath.Join(commitsDir, newCommit.ID)
	if err := os.WriteFile(commitFilePath, commitData, 0644); err != nil {
		log.Fatalf("Error writing commit file: %v", err)
	}

	// Update Metadata with the new commit
	metadataFile := ".gitx/metadata.json"
	if err := metadata_operations.UpdateMetadata(metadataFile, newCommit, "", ""); err != nil {
		log.Fatalf("Error updating metadata: %v", err)
	}

	// Update HEAD to point to the new commit
	if err := vcs_operations.UpdateHEAD(newCommit.ID); err != nil {
		log.Fatalf("Error updating HEAD: %v", err)
	}

	// Clear the staging area
	for k := range *stagingArea {
		delete(*stagingArea, k)
	}

	fmt.Printf("Commit created with ID: %s and message: %s\n", newCommit.ID, newCommit.Message)
}

func createInitialCommit() models.Commit {
	// Create an empty tree
	emptyTree := vcs_operations.CreateEmptyTree()

	// Set author and committer information
	author := vcs_operations.GetCurrentUser()
	committer := author

	// Set commit message and timestamp
	message := "Initial commit"
	timestamp := time.Now()

	// Generate commit ID
	commitID, err := vcs_operations.GenerateCommitID(emptyTree, nil, author, committer, timestamp)
	if err != nil {
		log.Fatalf("Error generating commit ID: %v", err)
	}

	// Create the initial commit object
	initialCommit := models.Commit{
		ID:        commitID,
		Parent:    nil, // No parent commit for the initial commit
		Tree:      emptyTree,
		Message:   message,
		Author:    author,
		Committer: committer,
		Timestamp: timestamp,
	}

	return initialCommit
}

// StatusHandler compares the files in the staging area with the tracked files in the metadata and the files in the working directory.
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

// getAllFilesInDir returns a list of all files in a directory.
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
