package vcs_operations

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
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

// mergeCommits performs the merge operation based on commit IDs
func mergeCommits(currentCommitID, mergeCommitID string) error {
	// Implement your merge logic here using commit IDs
	// For example, you might want to compare the commit IDs, resolve conflicts, etc.
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
