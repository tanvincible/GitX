package hash

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// SHA1Hash calculates the SHA-1 hash of the given file's content in Git blob format and stores the blob.
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
	header := fmt.Sprintf("blob %d\x00", info.Size())
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

	// Store the blob in the object database
	objectDir := filepath.Join(".gitx", "objects", hashedString[:2])
	objectFile := filepath.Join(objectDir, hashedString[2:])
	if err := os.MkdirAll(objectDir, os.ModePerm); err != nil {
		return "", err
	}

	// Reopen the file for reading
	file.Seek(0, 0)
	blobFile, err := os.Create(objectFile)
	if err != nil {
		return "", err
	}
	defer blobFile.Close()

	// Write the blob header and content to the object file
	if _, err := blobFile.Write([]byte(header)); err != nil {
		return "", err
	}
	if _, err := io.Copy(blobFile, file); err != nil {
		return "", err
	}

	return hashedString, nil
}
