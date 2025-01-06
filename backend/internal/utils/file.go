package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func SaveFile(file *multipart.File, header *multipart.FileHeader) (string, error) {
	// Create files directory if it doesn't exist
	err := os.MkdirAll("files", 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	// Create file path
	filePath := filepath.Join("files", header.Filename)

	// Create the file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer dst.Close()

	// Copy uploaded file data to destination file
	if _, err := io.Copy(dst, *file); err != nil {
		return "", fmt.Errorf("failed to copy file: %v", err)
	}

	return filePath, nil
}
