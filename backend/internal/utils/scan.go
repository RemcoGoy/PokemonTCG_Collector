package utils

import (
	"backend/internal/types"
	"encoding/binary"
	"encoding/csv"
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

func readCardHashes() ([]types.CardHash, error) {
	// Read card hashes from CSV file
	file, err := os.Open("data/card_hashes.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to open card hashes file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, err = reader.Read() // Skip header row
	if err != nil {
		return nil, fmt.Errorf("failed to read header row: %v", err)
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read records: %v", err)
	}

	cardHashes := make([]types.CardHash, 0)
	for _, record := range records {
		cardHashes = append(cardHashes, types.CardHash{
			ID:         record[0],
			Perceptual: record[1],
			Difference: record[2],
			Wavelet:    record[3],
			Color:      record[4],
			TCGID:      record[5],
		})
	}

	return cardHashes, nil
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

func FindClosestCard(hash *goimagehash.ExtImageHash) (types.CardHash, error) {
	cardHashes, err := readCardHashes()
	if err != nil {
		return types.CardHash{}, fmt.Errorf("failed to read card hashes: %v", err)
	}

	// Find closest match by calculating hamming distance
	minDistance := -1
	closestIdx := 0

	for i, cardHash := range cardHashes {
		// Convert stored hash string to ExtImageHash
		pHash := "p:" + cardHash.Perceptual
		storedHash, err := stringToExtHash(pHash)
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

	return cardHashes[closestIdx], nil
}

func PhashImage(file multipart.File, header *multipart.FileHeader) (*goimagehash.ExtImageHash, error) {
	img, err := fileToImage(file, header)
	if err != nil {
		return nil, fmt.Errorf("failed to convert file to image: %v", err)
	}

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
