package utils

import (
	"backend/internal/types"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"os"

	"github.com/corona10/goimagehash"
)

const (
	PHASH_SIZE       = 32
	PHASH_BLOCK_SIZE = 32
	extStrFmt        = "%1s:%s"
)

func readCardHashesGob() ([]types.CardHashGob, error) {
	file, err := os.Open("data/card_hashes.gob")
	if err != nil {
		return nil, fmt.Errorf("failed to open card hashes file: %v", err)
	}
	defer file.Close()

	var cardHashes [][]string
	err = gob.NewDecoder(file).Decode(&cardHashes)
	if err != nil {
		return nil, fmt.Errorf("failed to decode card hashes: %v", err)
	}

	cardHashesGob := make([]types.CardHashGob, 0)
	for _, cardHash := range cardHashes {
		cardHashesGob = append(cardHashesGob, types.CardHashGob{
			ID:   cardHash[0],
			Hash: cardHash[1],
		})
	}

	return cardHashesGob, nil
}

func stringToExtHash(s string) (*goimagehash.ExtImageHash, error) {
	var kindStr string
	var hashStr string
	_, err := fmt.Sscanf(s, extStrFmt, &kindStr, &hashStr)
	if err != nil {
		return nil, errors.New("Couldn't parse string " + s)
	}

	hexBytes, err := hex.DecodeString(hashStr)
	if err != nil {
		return nil, err
	}

	var hash []uint64
	lenOfByte := 8
	for i := 0; i < len(hexBytes)/lenOfByte; i++ {
		startIndex := i * lenOfByte
		endIndex := startIndex + lenOfByte
		hashUint64 := binary.BigEndian.Uint64(hexBytes[startIndex:endIndex])
		hash = append(hash, hashUint64)
	}

	kind := goimagehash.Unknown
	switch kindStr {
	case "a":
		kind = goimagehash.AHash
	case "p":
		kind = goimagehash.PHash
	case "d":
		kind = goimagehash.DHash
	case "w":
		kind = goimagehash.WHash
	}

	return goimagehash.NewExtImageHash(hash, kind, len(hash)*64), nil
}

func FindClosestCard(hash *goimagehash.ExtImageHash) (string, error) {
	// cardHashes, err := readCardHashes()
	cardHashes, err := readCardHashesGob()
	if err != nil {
		return "", fmt.Errorf("failed to read card hashes: %v", err)
	}

	// Find closest match by calculating hamming distance
	minDistance := -1
	closestIdx := 0

	for i, cardHash := range cardHashes {
		// Convert stored hash string to ExtImageHash
		storedHash, err := stringToExtHash(cardHash.Hash)
		if err != nil {
			fmt.Println("Error loading stored hash:", err)
			continue
		}

		distance, err := hash.Distance(storedHash)
		if err != nil {
			fmt.Println("Error calculating distance:", err)
			continue
		}

		// Update minimum distance if this is the first hash or if distance is smaller
		if minDistance == -1 || distance < minDistance {
			minDistance = distance
			closestIdx = i
		}
	}

	return cardHashes[closestIdx].ID, nil
}

func PhashFile(file multipart.File, header *multipart.FileHeader) (*goimagehash.ExtImageHash, error) {
	img, err := fileToImage(file, header)
	if err != nil {
		return nil, fmt.Errorf("failed to convert file to image: %v", err)
	}

	return PhashImage(img)
}

func PhashImage(img image.Image) (*goimagehash.ExtImageHash, error) {
	// Calculate perceptual hash
	hash, err := goimagehash.ExtPerceptionHash(img, PHASH_SIZE, PHASH_BLOCK_SIZE)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate phash: %v", err)
	}

	return hash, nil
}

func fileToImage(file multipart.File, header *multipart.FileHeader) (image.Image, error) {
	// Decode image
	var img image.Image
	var err error
	if header.Header.Get("Content-Type") == "image/png" {
		img, err = png.Decode(file)
	} else if header.Header.Get("Content-Type") == "image/jpeg" {
		img, err = jpeg.Decode(file)
	} else {
		return nil, fmt.Errorf("unsupported image format: %v", header.Header.Get("Content-Type"))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	return img, nil
}
