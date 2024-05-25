package vcs_operations

import (
	"fmt"
	"log"
	"os"
	"os/exec"
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
	cmd := exec.Command("git", "checkout", "-b", branchName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create branch: %v", err)
	}

	fmt.Printf("Created branch: %s\n", branchName)
	return nil
}

// ListBranches lists all Git branches.
func ListBranches() error {
	cmd := exec.Command("git", "branch")
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to list branches: %v", err)
	}
	fmt.Printf("Branches:\n%s\n", out)
	return nil
}

// SwitchBranch switches to the specified Git branch.
func SwitchBranch(branchName string) error {
	cmd := exec.Command("git", "checkout", branchName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to switch branch: %v", err)
	}

	fmt.Printf("Switched to branch: %s\n", branchName)
	return nil
}

// DeleteBranch deletes the specified Git branch.
func DeleteBranch(branchName string) error {
	cmd := exec.Command("git", "branch", "-d", branchName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete branch: %v", err)
	}

	fmt.Printf("Deleted branch: %s\n", branchName)
	return nil
}
