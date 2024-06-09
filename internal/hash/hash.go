package hash

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// SHA1Hash calculates the SHA-1 hash of the given file's content in Git blob format.
func SHA1Hash(filePath string) (string, error) {
	// Open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Get file info to retrieve its size
	info, err := file.Stat()
	if err != nil {
		return "", err
	}

	// Create a new SHA-1 hash instance
	h := sha1.New()

	// Write the blob header to the hash instance
	header := fmt.Sprintf("blob %d\\0", info.Size())
	if _, err := h.Write([]byte(header)); err != nil {
		return "", err
	}

	// Write the file content to the hash instance
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}

	// Calculate the hash
	hashed := h.Sum(nil)

	// Encode the hashed result to hexadecimal string
	hashedString := hex.EncodeToString(hashed)

	return hashedString, nil
}
