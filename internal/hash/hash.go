package hash

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
)

// SHA1Hash calculates the SHA-1 hash of the given content.
func SHA1Hash(content string) (string, error) {
	// Ensure content is not empty
	if content == "" {
		return "", errors.New("empty content provided")
	}

	// Create a new SHA-1 hash instance
	h := sha1.New()

	// Write content to the hash instance
	_, err := h.Write([]byte(content))
	if err != nil {
		return "", err
	}

	// Calculate the hash
	hashed := h.Sum(nil)

	// Encode the hashed result to hexadecimal string
	hashedString := hex.EncodeToString(hashed)

	return hashedString, nil
}
