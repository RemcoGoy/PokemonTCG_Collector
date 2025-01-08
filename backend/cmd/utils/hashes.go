package main

import (
	"backend/internal/utils"
	"encoding/csv"
	"fmt"
	"image"
	"log"
	"net/http"
	"os"
	"time"

	tcg "github.com/PokemonTCG/pokemon-tcg-sdk-go-v2/pkg"
	"github.com/PokemonTCG/pokemon-tcg-sdk-go-v2/pkg/request"
	"github.com/joho/godotenv"
)

func createCsvWriter(filename string) (*csv.Writer, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return csv.NewWriter(file), nil
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
func main() {
	start := time.Now()
	log.Println("Starting hash generation process...")

	if err := os.Remove("./data/hashes.csv"); err != nil && !os.IsNotExist(err) {
		log.Fatal("Error removing existing hashes.csv file:", err)
	}
	log.Println("Cleaned up any existing hashes.csv file")

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	API_KEY := os.Getenv("TCG_API_KEY")
	tcg_client := tcg.NewClient(API_KEY)
	log.Println("Successfully loaded API key and created TCG client")

	writer, err := createCsvWriter("./data/hashes.csv")
	if err != nil {
		log.Fatal("Error creating CSV writer:", err)
	}
	defer writer.Flush()
	log.Println("Successfully created CSV writer")

	numWorkers := 4
	numPages := 10
	jobs := make(chan int, numPages)
	results := make(chan [][]string, numPages)

	log.Printf("Starting %d workers to process %d pages...", numWorkers, numPages)

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

	// Send jobs
	log.Println("Sending jobs to workers...")
	for page := 1; page <= numPages; page++ {
		jobs <- page
	}
	close(jobs)
	log.Println("All jobs sent to workers")

	// Collect results
	log.Println("Collecting results...")
	var hashes [][]string
	for i := 0; i < numPages; i++ {
		pageHashes := <-results
		hashes = append(hashes, pageHashes...)
		log.Printf("Collected results for page %d", i+1)
	}

	log.Printf("Writing %d hash records to CSV...", len(hashes))
	for _, hash := range hashes {
		writeRecord(writer, hash)
	}

	elapsed := time.Since(start)
	log.Printf("Processing completed. Total time: %s", elapsed)
}
