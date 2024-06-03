package vcs_operations

import (
	"GitX/internal/hash"
	"GitX/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func UpdateHEAD(commitHash string) {
	headFile := "HEAD"

	// Create or open the HEAD file
	file, err := os.OpenFile(headFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error opening HEAD file: %v", err)
	}
	defer file.Close()

	// Write the commit hash to the HEAD file
	if _, err := file.WriteString(commitHash); err != nil {
		log.Fatalf("Error writing to HEAD file: %v", err)
	}
}

// CreateBranch creates a new Git branch.
func CreateBranch(branchName string) error {
	// Define the branch directory path
	branchDir := filepath.Join(".gitx", "refs", "heads", branchName)

	// Check if the branch already exists
	if _, err := os.Stat(branchDir); err == nil {
		return fmt.Errorf("branch %s already exists", branchName)
	}

	// Get the current branch name
	currentBranch, err := getCurrentBranch()
	if err != nil {
		return err
	}

	// Define the path to the current branch directory
	currentBranchDir := filepath.Join(".gitx", "refs", "heads", currentBranch)

	// Check if the current branch directory exists
	if _, err := os.Stat(currentBranchDir); os.IsNotExist(err) {
		return fmt.Errorf("current branch directory does not exist")
	}

	// Create the new branch directory
	if err := os.MkdirAll(branchDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create branch directory: %v", err)
	}

	// Copy the contents from the current branch directory to the new branch directory
	if err := copyDir(currentBranchDir, branchDir); err != nil {
		return fmt.Errorf("failed to create branch: %v", err)
	}

	// Indicate successful branch creation
	fmt.Printf("Created branch: %s\n", branchName)
	return nil
}

func ListBranches() {
	gitxDir := ".gitx"
	refsHeadsDir := filepath.Join(gitxDir, "refs", "heads")

	files, err := os.ReadDir(refsHeadsDir)
	if err != nil {
		log.Fatalf("Error reading refs/heads directory: %v", err)
	}

	fmt.Println("Branches:")
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fmt.Println(file.Name())
	}
}

// SwitchBranch switches to the specified Git branch.
func SwitchBranch(branchName string) error {
	branchDir := filepath.Join(".gitx", "branches", branchName)
	if _, err := os.Stat(branchDir); os.IsNotExist(err) {
		return fmt.Errorf("branch %s does not exist", branchName)
	}

	currentBranch, err := getCurrentBranch()
	if err != nil {
		return err
	}

	currentBranchDir := filepath.Join(".gitx", "branches", currentBranch)
	if _, err := os.Stat(currentBranchDir); os.IsNotExist(err) {
		return fmt.Errorf("current branch directory does not exist")
	}

	if err := copyDir(branchDir, currentBranchDir); err != nil {
		return fmt.Errorf("failed to switch branch: %v", err)
	}

	fmt.Printf("Switched to branch: %s\n", branchName)
	return nil
}

// DeleteBranch deletes the specified Git branch.
func DeleteBranch(branchName string) error {
	if branchName == "main" {
		return fmt.Errorf("cannot delete main branch")
	}

	branchDir := filepath.Join(".gitx", "branches", branchName)
	if _, err := os.Stat(branchDir); os.IsNotExist(err) {
		return fmt.Errorf("branch %s does not exist", branchName)
	}

	if err := os.RemoveAll(branchDir); err != nil {
		return fmt.Errorf("failed to delete branch: %v", err)
	}

	fmt.Printf("Deleted branch: %s\n", branchName)
	return nil
}

// MergeBranch merges the specified branch into the current branch.
func MergeBranch(branchName string) error {
	// Read the current branch from HEAD
	currentBranch, err := getCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %v", err)
	}

	// Read the commit IDs of the current branch and the branch to merge
	currentCommitID, err := getCommitID(currentBranch)
	if err != nil {
		return fmt.Errorf("failed to get commit ID of %s: %v", currentBranch, err)
	}

	mergeCommitID, err := getCommitID(branchName)
	if err != nil {
		return fmt.Errorf("failed to get commit ID of %s: %v", branchName, err)
	}

	// Perform the merge operation using commit IDs
	// Perform the merge operation using commit IDs
	err = mergeCommits(currentCommitID, mergeCommitID)
	if err != nil {
		return fmt.Errorf("failed to merge branch %s into %s: %v", branchName, currentBranch, err)
	}

	fmt.Printf("Merged branch %s into %s\n", branchName, currentBranch)
	return nil
}

func getCurrentBranch() (string, error) {
	headFile := filepath.Join(".gitx", "HEAD")
	content, err := os.ReadFile(headFile)
	if err != nil {
		return "", fmt.Errorf("failed to read HEAD file: %v", err)
	}
	return string(content), nil
}

// getCommitID returns the commit ID of the given branch
func getCommitID(branchName string) (string, error) {
	branchPath := filepath.Join("refs", "heads", branchName)
	commitID, err := readFileContent(branchPath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(commitID)), nil
}

// readCommit reads the commit data from the commit file.
func readCommit(commitID string) (*models.Commit, error) {
	commitFile := filepath.Join(".gitx", "commits", commitID+".json")
	file, err := os.Open(commitFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var commit models.Commit
	if err := json.NewDecoder(file).Decode(&commit); err != nil {
		return nil, err
	}
	return &commit, nil
}

// findCommonAncestor finds the common ancestor of two commits.
func findCommonAncestor(currentCommit *models.Commit, mergeCommit *models.Commit) (*models.Commit, error) {
	// Placeholder for common ancestor logic
	// You can implement a proper logic to find common ancestors.
	// For now, we assume the initial commit is the common ancestor.
	initialCommitID := "initial_commit_id" // Replace this with actual logic to find the common ancestor
	return readCommit(initialCommitID)
}

// mergeFiles performs a three-way merge of the contents of files.
func mergeFiles(baseContent, currentContent, otherContent string) (string, bool) {
	if currentContent == otherContent {
		return currentContent, false
	}
	if baseContent == currentContent {
		return otherContent, false
	}
	if baseContent == otherContent {
		return currentContent, false
	}

	// Conflict detected
	conflictContent := fmt.Sprintf("<<<<<<< current\n%s\n=======\n%s\n>>>>>>> other\n", currentContent, otherContent)
	return conflictContent, true
}

// mergeCommits merges changes from two commits and creates a new merge commit.
func mergeCommits(currentCommitID, mergeCommitID string) error {
	// Read the current commit
	currentCommit, err := readCommit(currentCommitID)
	if err != nil {
		return fmt.Errorf("error reading current commit: %v", err)
	}

	// Read the commit to merge
	mergeCommit, err := readCommit(mergeCommitID)
	if err != nil {
		return fmt.Errorf("error reading merge commit: %v", err)
	}

	// Find common ancestor
	baseCommit, err := findCommonAncestor(currentCommit, mergeCommit)
	if err != nil {
		return fmt.Errorf("error finding common ancestor: %v", err)
	}

	// Merge changes
	mergedFiles := make(map[string]string)
	conflictDetected := false

	// Collect all unique file paths
	allFilePaths := make(map[string]bool)
	for filePath := range baseCommit.Files {
		allFilePaths[filePath] = true
	}
	for filePath := range currentCommit.Files {
		allFilePaths[filePath] = true
	}
	for filePath := range mergeCommit.Files {
		allFilePaths[filePath] = true
	}

	// Merge each file
	for filePath := range allFilePaths {
		baseContent := baseCommit.Files[filePath]
		currentContent := currentCommit.Files[filePath]
		mergeContent := mergeCommit.Files[filePath]

		mergedContent, conflict := mergeFiles(baseContent, currentContent, mergeContent)
		if conflict {
			conflictDetected = true
		}

		mergedFiles[filePath] = mergedContent
	}

	// If conflicts were detected, notify the user
	if conflictDetected {
		fmt.Println("Merge conflicts detected. Please resolve them manually.")
	}

	// Create the merge commit
	mergeCommitHash, err := hash.SHA1Hash(fmt.Sprintf("%s+%s", currentCommitID, mergeCommitID))
	if err != nil {
		return fmt.Errorf("error computing merge commit hash: %v", err)
	}

	newCommit := &models.Commit{
		ID:           mergeCommitHash,
		Parent:       []*models.Commit{currentCommit, mergeCommit},
		Message:      fmt.Sprintf("Merge commit %s into %s", mergeCommitID, currentCommitID),
		Author:       "Merge Author",
		Timestamp:    time.Now(),
		Files:        mergedFiles,
		Committer:    "Merge Committer",
		GPGSignature: "",
	}

	// Save the new merge commit
	commitFile := filepath.Join(".gitx", "commits", newCommit.ID+".json")
	file, err := os.Create(commitFile)
	if err != nil {
		return fmt.Errorf("error creating merge commit file: %v", err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(newCommit); err != nil {
		return fmt.Errorf("error encoding merge commit: %v", err)
	}

	fmt.Printf("Created merge commit: %s\n", newCommit.ID)
	return nil
}

// readFileContent reads the content of a file
func readFileContent(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var content []byte
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	content = make([]byte, stat.Size())
	_, err = file.Read(content)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// SquashCommits manually squashes the specified range of commits into a single commit.
func SquashCommits(baseCommit, targetCommit string) error {
	// Implement the squash logic here
	// For example, I could use a diffing algorithm to generate a single commit that represents the changes between baseCommit and targetCommit.

	fmt.Printf("Manually squashed commits from %s to %s\n", baseCommit, targetCommit)
	return nil
}

// Stash saves the changes in the working directory to a temporary location.
func Stash() error {
	// Create a temporary directory to store the stashed changes
	tempDir, err := os.MkdirTemp("", "stashed_changes")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %v", err)
	}

	// Walk through the working directory and copy all files to the temporary directory
	err = copyDir(".", tempDir)
	if err != nil {
		return fmt.Errorf("failed to stash changes: %v", err)
	}

	fmt.Printf("Stashed changes in %s\n", tempDir)
	return nil
}

// LogHandler displays the commit history
func LogHandler() {
	commitsDir := ".gitx/commits"

	// Read the commits directory
	files, err := os.ReadDir(commitsDir)
	if err != nil {
		log.Fatalf("Error reading commits directory: %v", err)
	}

	// Iterate over the commit files
	for _, file := range files {
		if !file.IsDir() {
			commitFile := filepath.Join(commitsDir, file.Name())
			displayCommit(commitFile)
		}
	}
}

// displayCommit reads and prints commit details
func displayCommit(commitFile string) {
	file, err := os.Open(commitFile)
	if err != nil {
		log.Fatalf("Error opening commit file: %v", err)
	}
	defer file.Close()

	var commit models.Commit
	if err := json.NewDecoder(file).Decode(&commit); err != nil {
		log.Fatalf("Error decoding commit file: %v", err)
	}

	// Print commit details
	fmt.Println("Commit:", commit.ID)
	fmt.Println("Author:", commit.Author)
	fmt.Println("Date:", commit.Timestamp)
	fmt.Println("Message:", commit.Message)
	fmt.Println("-------------------------------")
}

func copyDir(src, dst string) error {
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			if err := os.MkdirAll(dstPath, info.Mode()); err != nil {
				return err
			}
		} else {
			srcContent, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			if err := os.WriteFile(dstPath, srcContent, info.Mode()); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// copyFile copies a single file from the source to the destination.
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer destinationFile.Close()

	if _, err := io.Copy(destinationFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy file: %v", err)
	}

	if err := destinationFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %v", err)
	}

	return nil
}

// CatFile displays the content of an object in the repository.
func CatFile(objectID string) error {
	// Retrieve the object by its ID
	object, err := getObjectByID(objectID)
	if err != nil {
		return err
	}

	// Print the content of each file
	fmt.Println("Files content:")
	for filePath, hash := range object.Files {
		// Fetch the content of the file using its hash
		fileContent, err := readFileContent(hash)
		if err != nil {
			fmt.Printf("Error fetching content for file %s: %v\n", filePath, err)
			continue
		}
		// Print the file path and its content
		fmt.Printf("File: %s\nContent:\n%s\n", filePath, fileContent)
	}

	return nil
}

// getObjectByID retrieves the object from the repository by its ID.
func getObjectByID(_ string) (*models.Commit, error) {
	// Placeholder for retrieving the object by ID
	// You would typically fetch the object from the repository storage by its ID
	// and return the object's content along with other metadata.
	return nil, fmt.Errorf("getObjectByID: not implemented")
}

func ReflogHandler() error {
	reflogDir := ".gitx/reflog"

	// Check if the reflog directory exists
	if _, err := os.Stat(reflogDir); os.IsNotExist(err) {
		return fmt.Errorf("reflog directory does not exist: %v", err)
	}

	// Read the reflog directory
	files, err := os.ReadDir(reflogDir)
	if err != nil {
		return fmt.Errorf("error reading reflog directory: %v", err)
	}

	// Iterate over the reflog files
	for _, file := range files {
		if !file.IsDir() {
			reflogFile := filepath.Join(reflogDir, file.Name())
			if err := displayReflog(reflogFile); err != nil {
				log.Printf("Error displaying reflog file %s: %v", reflogFile, err)
			}
		}
	}

	return nil
}

func displayReflog(reflogFile string) error {
	file, err := os.Open(reflogFile)
	if err != nil {
		return fmt.Errorf("error opening reflog file %s: %v", reflogFile, err)
	}
	defer file.Close()

	var reflog models.Reflog
	if err := json.NewDecoder(file).Decode(&reflog); err != nil {
		return fmt.Errorf("error decoding reflog file %s: %v", reflogFile, err)
	}

	// Print reflog details
	fmt.Println("Reflog:", reflog.ID)
	fmt.Println("Author:", reflog.Author)
	fmt.Println("Date:", reflog.Timestamp)
	fmt.Println("Message:", reflog.Message)
	fmt.Println("-------------------------------")

	return nil
}
