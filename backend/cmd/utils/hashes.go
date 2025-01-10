package main

import (
	"backend/internal/utils"
	"encoding/csv"
	"encoding/gob"
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

const (
	PAGE_SIZE  = 10 // Number of cards per page
	WORKERS    = 20 // Number of workers
	BATCH_SIZE = 20 // Number of pages per batch before writing to CSV
)

var (
	TOTAL_CARDS = -1 // Total number of cards to process
	TOTAL_PAGES = -1 // Total number of pages to process
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

func verifyGob() {
	gobFile, err := os.Open("./data/hashes.gob")
	if err != nil {
		log.Fatal("Error opening gob file:", err)
	}
	defer gobFile.Close()

	decoder := gob.NewDecoder(gobFile)
	var hashes [][]string
	if err := decoder.Decode(&hashes); err != nil {
		log.Fatal("Error decoding gob file:", err)
	}

	fmt.Println("\nFirst 5 entries:")
	for i := 0; i < 5 && i < len(hashes); i++ {
		fmt.Printf("%v\n", hashes[i])
	}

	fmt.Printf("\nTotal number of entries: %d\n", len(hashes))
}

func main() {
	cleanupFlag := flag.Bool("cleanup", true, "Whether to cleanup existing hashes.csv file")
	formatFlag := flag.String("format", "csv", "Output format (csv or gob)")
	flag.Parse()

	start := time.Now()
	log.Println("Starting hash generation process...")

	TOTAL_CARDS = getTotalCards()
	log.Printf("Total cards to process: %d", TOTAL_CARDS)

	if *cleanupFlag {
		cleanup()
		log.Println("Cleaned up any existing hashes.csv file")
	}

	tcg_client := createClient()
	log.Println("Successfully loaded API key and created TCG client")

	var writer *csv.Writer
	var gobFile *os.File
	var err error

	if *formatFlag == "csv" {
		writer, err = createCsvWriter("./data/hashes.csv")
		if err != nil {
			log.Fatal("Error creating CSV writer:", err)
		}
		defer writer.Flush()
		log.Println("Successfully created CSV writer")
		writeRecord(writer, []string{"id", "hash"})
	} else if *formatFlag == "gob" {
		gobFile, err = os.Create("./data/hashes.gob")
		if err != nil {
			log.Fatal("Error creating gob file:", err)
		}
		defer gobFile.Close()
		log.Println("Successfully created gob file")
	} else {
		log.Fatal("Invalid format specified. Must be 'csv' or 'gob'")
	}

	TOTAL_PAGES = (TOTAL_CARDS + 9) / PAGE_SIZE

	jobs := make(chan int, BATCH_SIZE)
	results := make(chan [][]string, BATCH_SIZE)

	// Start workers
	for w := 0; w < WORKERS; w++ {
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

	log.Printf("Processing %d total pages in batches of %d...", TOTAL_PAGES, BATCH_SIZE)

	var allHashes [][]string

	// Process in batches
	for batchStart := 1; batchStart <= TOTAL_PAGES; batchStart += BATCH_SIZE {
		batchEnd := batchStart + BATCH_SIZE - 1
		if batchEnd > TOTAL_PAGES {
			batchEnd = TOTAL_PAGES
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
		}

		if *formatFlag == "csv" {
			// Write batch results to CSV
			for _, hash := range batchHashes {
				writeRecord(writer, hash)
			}
		} else {
			// Collect all hashes for gob encoding
			allHashes = append(allHashes, batchHashes...)
		}

		log.Printf("Completed batch %d to %d", batchStart, batchEnd)
	}

	close(jobs)

	if *formatFlag == "gob" {
		encoder := gob.NewEncoder(gobFile)
		if err := encoder.Encode(allHashes); err != nil {
			log.Fatal("Error encoding to gob:", err)
		}
		log.Println("Successfully wrote all hashes to gob file")

		verifyGob()
	}

	elapsed := time.Since(start)
	log.Printf("Processing completed. Total time: %s", elapsed)
}
