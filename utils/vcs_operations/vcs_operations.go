package vcs_operations

import (
	"log"
	"os"
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
