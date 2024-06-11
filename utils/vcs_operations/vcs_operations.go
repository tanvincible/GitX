package vcs_operations

import (
	"GitX/internal/hash"
	"GitX/models"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// UpdateHEAD updates the HEAD file with the reference to the latest commit on the current branch.
func UpdateHEAD(commitHash string) error {
	headFile := ".gitx/HEAD"

	// Write the commit hash to the HEAD file
	if err := os.WriteFile(headFile, []byte(commitHash), 0644); err != nil {
		return fmt.Errorf("error writing to HEAD file: %w", err)
	}

	return nil
}

// GetCurrentHeadCommit retrieves the current commit that HEAD is pointing to.
func GetCurrentHeadCommit() string {
	headFile := ".gitx/HEAD"
	content, err := os.ReadFile(headFile)
	if err != nil {
		if os.IsNotExist(err) {
			// If the HEAD file doesn't exist, it means there are no commits yet.
			return ""
		}
		log.Fatalf("Error reading HEAD file: %v", err)
	}

	headContent := string(content)
	// Check if HEAD is a reference to a branch
	if strings.HasPrefix(headContent, "refs/heads/") {
		// Extract the branch name
		branchName := strings.TrimPrefix(headContent, "refs/heads/")
		// Check if the branch has any commits
		branchCommitsFile := fmt.Sprintf(".gitx/refs/heads/%s", branchName)
		if _, err := os.Stat(branchCommitsFile); os.IsNotExist(err) {
			// If the branch commits file doesn't exist, there are no commits on this branch
			return ""
		} else {
			// Read the last commit hash from the branch commits file
			lastCommitHash, err := os.ReadFile(branchCommitsFile)
			if err != nil {
				log.Fatalf("Error reading branch commits file: %v", err)
			}
			return string(lastCommitHash)
		}
	}

	// If HEAD is not a branch reference, assume it's a commit hash
	return headContent
}

// CreateTreeFromStagingArea takes a map of file paths to their hashes and creates a Tree object.
func CreateTreeFromStagingArea(stagingArea map[string]string) (*models.Tree, error) {
	tree := &models.Tree{
		Entries: []models.TreeEntry{},
	}

	// Sort the file paths to ensure consistent tree hash
	var paths []string
	for path := range stagingArea {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	// Create TreeEntries for each file in the staging area
	for _, path := range paths {
		hashValue := stagingArea[path]
		entry := models.TreeEntry{
			Name: path,
			Mode: "100644", // This mode represents a regular non-executable file
			ID:   hashValue,
			Type: "blob", // Assuming all entries are files for simplicity
		}
		tree.Entries = append(tree.Entries, entry)
	}

	// Generate the SHA-1 hash for the tree
	hash := sha1.New()
	for _, entry := range tree.Entries {
		// Create a string representation of the entry for hashing
		entryStr := fmt.Sprintf("%s %s %s\t%s", entry.Mode, entry.Type, entry.ID, entry.Name)
		hash.Write([]byte(entryStr))
	}
	tree.ID = hex.EncodeToString(hash.Sum(nil))

	return tree, nil
}

// GetCommitByHash retrieves a commit object by its hash.
func GetCommitByHash(commitHash string) (*models.Commit, error) {
	commitsDir := ".gitx/commits"
	commitFilePath := filepath.Join(commitsDir, commitHash)

	// Read the commit file using os.ReadFile
	commitData, err := os.ReadFile(commitFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("commit %s does not exist", commitHash)
		}
		return nil, fmt.Errorf("error reading commit file: %w", err)
	}

	// Unmarshal the commit data into a Commit object
	var commit models.Commit
	if err := json.Unmarshal(commitData, &commit); err != nil {
		return nil, fmt.Errorf("error unmarshaling commit data: %w", err)
	}

	return &commit, nil
}

// GetCurrentUser retrieves the current user from the system.
func GetCurrentUser() string {
	// Placeholder for getting the current user
	return "John Doe"
}

func CreateEmptyTree() *models.Tree {
    return &models.Tree{
        Entries: []models.TreeEntry{},
    }
}

// Function to check if a branch exists by looking for its reference file
func branchExists(branchName string) bool {
	refsHeadsDir := filepath.Join(".gitx", "refs", "heads")
	branchRefPath := filepath.Join(refsHeadsDir, branchName)
	if _, err := os.Stat(branchRefPath); err == nil {
		return true
	}
	return false
}

// getCurrentBranch reads the current branch from the HEAD file.
func getCurrentBranch() (string, error) {
	headFile := filepath.Join(".gitx", "HEAD")
	content, err := os.ReadFile(headFile)
	if err != nil {
		return "", fmt.Errorf("failed to read HEAD file: %v", err)
	}

	// Parse the content to extract the branch name
	refPrefix := "ref: refs/heads/"
	if strings.HasPrefix(string(content), refPrefix) {
		// The branch name is the part after the prefix
		return strings.TrimSpace(strings.TrimPrefix(string(content), refPrefix)), nil
	}

	return "", fmt.Errorf("HEAD file does not contain a valid branch reference")
}

// GenerateCommitID generates a commit ID based on the tree hash, parent commit IDs, and other commit information.
func GenerateCommitID(tree *models.Tree, parents []*models.Commit, message, author string, timestamp time.Time) (string, error) {
	// Create a new SHA-1 hash instance
	h := sha1.New()

	// Serialize the tree hash
	treeHash := tree.ID // Use the ID field of the tree object
	h.Write([]byte(fmt.Sprintf("tree %s\n", treeHash)))

	// Serialize parent commits
	for _, parent := range parents {
		h.Write([]byte(fmt.Sprintf("parent %s\n", parent.ID)))
	}

	// Serialize author information
	authorInfo := fmt.Sprintf("author %s %d +0000\n", author, timestamp.Unix())
	h.Write([]byte(authorInfo))

	// Serialize committer information (same as author for simplicity)
	committerInfo := fmt.Sprintf("committer %s %d +0000\n", author, timestamp.Unix())
	h.Write([]byte(committerInfo))

	// Serialize commit message
	h.Write([]byte(fmt.Sprintf("\n%s\n", message)))

	// Calculate the hash
	hashed := h.Sum(nil)

	// Encode the hashed result to hexadecimal string
	hashedString := hex.EncodeToString(hashed)

	return hashedString, nil
}

// CreateBranch creates a new Git branch.
func CreateBranch(branchName string) error {
	if branchExists(branchName) {
		return fmt.Errorf("branch '%s' already exists", branchName)
	}
	gitxDir := ".gitx"
	refsHeadsDir := filepath.Join(gitxDir, "refs", "heads")
	branchRefPath := filepath.Join(refsHeadsDir, branchName)

	// Check if the branch already exists
	if _, err := os.Stat(branchRefPath); err == nil {
		return fmt.Errorf("branch '%s' already exists", branchName)
	}

	// Initialize the new branch ref file with an empty content or placeholder
	initialContent := []byte("") // Empty content since there are no commits yet

	// Write the initial content to the new branch ref file
	if err := os.WriteFile(branchRefPath, initialContent, 0644); err != nil {
		return fmt.Errorf("error initializing branch ref file: %v", err)
	}

	fmt.Printf("Branch '%s' created successfully.\n", branchName)
	return nil
}

// ListBranches lists all the Git branches in the repository.
func ListBranches() {
	gitxDir := ".gitx"
	refsHeadsDir := filepath.Join(gitxDir, "refs", "heads")

	// Read the current branch reference from the HEAD file
	headFile := filepath.Join(gitxDir, "HEAD")
	currentBranchRefBytes, err := os.ReadFile(headFile)
	if err != nil {
		log.Fatalf("Error reading HEAD file: %v", err)
	}
	currentBranchRef := strings.TrimSpace(string(currentBranchRefBytes))

	// Extract the branch name from the reference
	// Assuming the reference is in the format "refs/heads/branch-name"
	parts := strings.Split(currentBranchRef, "/")
	if len(parts) < 3 {
		log.Fatalf("Invalid branch reference in HEAD file")
	}
	currentBranch := parts[2] // The branch name is the third part of the reference

	files, err := os.ReadDir(refsHeadsDir)
	if err != nil {
		log.Fatalf("Error reading refs/heads directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		branchName := file.Name()
		if branchName == currentBranch {
			// Print the current branch in green with an asterisk
			fmt.Printf("\033[32m* %s\033[0m\n", branchName)
		} else {
			fmt.Println(branchName)
		}
	}
}

// SwitchBranch switches to the specified Git branch.
func SwitchBranch(branchName string) error {
	if !branchExists(branchName) {
		return fmt.Errorf("branch %s does not exist", branchName)
	}

	// Get the name of the current branch
	currentBranch, err := getCurrentBranch()
	if err != nil {
		return err
	}

	// Check if the target branch is the same as the current branch
	if currentBranch == branchName {
		fmt.Printf("Already on '%s'\n", branchName)
		return nil
	}

	// Path to the HEAD file
	headPath := filepath.Join(".gitx", "HEAD")

	// Content to write to the HEAD file, pointing to the new branch
	newHeadContent := []byte("ref: refs/heads/" + branchName + "\n")

	// Write the new HEAD content to point to the new branch
	if err := os.WriteFile(headPath, newHeadContent, 0644); err != nil {
		return fmt.Errorf("failed to switch branch: %v", err)
	}
	fmt.Println("Switched to branch:", branchName)
	return nil
}

// DeleteBranch deletes the specified Git branch.
func DeleteBranch(branchName string) error {
	if branchName == "" {
		return fmt.Errorf("branch name cannot be empty")
	}
	if !branchExists(branchName) {
		return fmt.Errorf("branch '%s' does not exist", branchName)
	}

	// Prevent deletion of the current branch
	currentBranch, err := getCurrentBranch()
	if err != nil {
		return err
	}
	if currentBranch == branchName {
		return fmt.Errorf("cannot delete the current branch")
	}

	// Path to the branch reference file
	refsHeadsDir := filepath.Join(".gitx", "refs", "heads")
	branchRefPath := filepath.Join(refsHeadsDir, branchName)

	// Delete the branch reference file
	if err := os.Remove(branchRefPath); err != nil {
		return fmt.Errorf("failed to delete branch: %v", err)
	}

	fmt.Printf("Branch '%s' deleted successfully.\n", branchName)
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

// copyDir copies the contents of a directory to another directory.
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

// ReflogHandler displays the reflog history.
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

// displayReflog reads and prints reflog details
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
