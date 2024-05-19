package compression

import (
	"bytes"
	"compress/zlib"
	"os"
)

// CompressFile compresses the content of a file located at filePath using zlib compression algorithm.
func CompressFile(filePath string) ([]byte, error) {
	// Read the content of the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Compress the content using zlib
	var compressed bytes.Buffer
	w := zlib.NewWriter(&compressed)
	_, err = w.Write(content)
	if err != nil {
		return nil, err
	}
	w.Close()

	return compressed.Bytes(), nil
}
