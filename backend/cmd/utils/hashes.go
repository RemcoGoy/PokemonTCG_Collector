package main

import (
	"backend/internal/utils"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	tcg "github.com/PokemonTCG/pokemon-tcg-sdk-go-v2/pkg"
	"github.com/PokemonTCG/pokemon-tcg-sdk-go-v2/pkg/request"
	"github.com/joho/godotenv"
)

func createCsvWriter(filename string) (*csv.Writer, error) {
	// Check if file exists first
	if _, err := os.Stat(filename); err == nil {
		// File exists, open it in append mode
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		return csv.NewWriter(file), nil
	} else {
		// File doesn't exist, create it
		file, err := os.Create(filename)
		if err != nil {
			return nil, err
		}

		return csv.NewWriter(file), nil
	}
}

func writeRecord(writer *csv.Writer, record []string) {
	err := writer.Write(record)

	if err != nil {
		fmt.Println("Error writing record to CSV:", err)
	}
}

func hashCards(page int, client tcg.Client) [][]string {
	cards, err := client.GetCards(request.OrderBy("set.releaseDate,number"), request.PageSize(10), request.Page(page))
	if err != nil {
		fmt.Println(err)
	}

	hashes := [][]string{}

	for _, card := range cards {
		resp, err := http.Get(card.Images.Large)
		if err != nil {
			fmt.Printf("Error downloading image: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		img, _, err := image.Decode(resp.Body)
		if err != nil {
			fmt.Printf("Error decoding image: %v\n", err)
			continue
		}

		hash, err := utils.PhashImage(img)
		if err != nil {
			fmt.Printf("Error calculating hash: %v\n", err)
			continue
		}

		hashes = append(hashes, []string{card.ID, hash.ToString()})
	}

	return hashes
}

func getTotalCards() int {
	url := "https://api.pokemontcg.io/v2/cards?pageSize=1&page=1"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error fetching cards:", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
	}
	var response struct {
		TotalCount int `json:"totalCount"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatal("Error parsing response:", err)
	}
	return response.TotalCount
}

func cleanup() {
	if err := os.Remove("./data/hashes.csv"); err != nil && !os.IsNotExist(err) {
		log.Fatal("Error removing existing hashes.csv file:", err)
	}
}

func createClient() tcg.Client {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	API_KEY := os.Getenv("TCG_API_KEY")
	return tcg.NewClient(API_KEY)
}

func main() {
	cleanupFlag := flag.Bool("cleanup", true, "Whether to cleanup existing hashes.csv file")
	flag.Parse()

	start := time.Now()
	log.Println("Starting hash generation process...")

	totalCards := getTotalCards()
	log.Printf("Total cards to process: %d", totalCards)

	if *cleanupFlag {
		cleanup()
		log.Println("Cleaned up any existing hashes.csv file")
	}

	tcg_client := createClient()
	log.Println("Successfully loaded API key and created TCG client")

	writer, err := createCsvWriter("./data/hashes.csv")
	if err != nil {
		log.Fatal("Error creating CSV writer:", err)
	}
	defer writer.Flush()
	log.Println("Successfully created CSV writer")

	numWorkers := 10
	batchSize := 10
	totalPages := (totalCards + 9) / 10

	jobs := make(chan int, batchSize)
	results := make(chan [][]string, batchSize)

	// Start workers
	for w := 0; w < numWorkers; w++ {
		go func(workerId int) {
			log.Printf("Worker %d started", workerId)
			for page := range jobs {
				log.Printf("Worker %d processing page %d", workerId, page)
				hashes := hashCards(page, tcg_client)
				results <- hashes
				log.Printf("Worker %d completed page %d", workerId, page)
			}
		}(w)
	}

	log.Printf("Processing %d total pages in batches of %d...", totalPages, batchSize)

	// Process in batches
	for batchStart := 1; batchStart <= totalPages; batchStart += batchSize {
		batchEnd := batchStart + batchSize - 1
		if batchEnd > totalPages {
			batchEnd = totalPages
		}

		currentBatchSize := batchEnd - batchStart + 1
		log.Printf("Processing batch from page %d to %d...", batchStart, batchEnd)

		// Send jobs for this batch
		for page := batchStart; page <= batchEnd; page++ {
			jobs <- page
		}

		// Collect results for this batch
		var batchHashes [][]string
		for i := 0; i < currentBatchSize; i++ {
			pageHashes := <-results
			batchHashes = append(batchHashes, pageHashes...)
			log.Printf("Collected results for page %d", batchStart+i)
		}

		// Write batch results to CSV
		log.Printf("Writing %d hash records to CSV...", len(batchHashes))
		for _, hash := range batchHashes {
			writeRecord(writer, hash)
		}

		log.Printf("Completed batch %d to %d", batchStart, batchEnd)
	}

	close(jobs)

	elapsed := time.Since(start)
	log.Printf("Processing completed. Total time: %s", elapsed)
}
