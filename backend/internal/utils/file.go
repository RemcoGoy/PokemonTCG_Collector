package utils

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/corona10/goimagehash"
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

func PhashImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	// Decode image
	var img image.Image
	var err error
	if header.Header.Get("Content-Type") == "image/png" {
		img, err = png.Decode(file)
	} else if header.Header.Get("Content-Type") == "image/jpeg" {
		img, err = jpeg.Decode(file)
	} else {
		return "", fmt.Errorf("unsupported image format: %v", header.Header.Get("Content-Type"))
	}
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %v", err)
	}

	// Calculate perceptual hash
	hash, err := goimagehash.PerceptionHash(img)
	if err != nil {
		return "", fmt.Errorf("failed to calculate phash: %v", err)
	}

	return hash.ToString(), nil
}
