package main

import (
	"backend/internal/utils"
	"encoding/csv"
	"fmt"
	"image"
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

	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}
	API_KEY := os.Getenv("TCG_API_KEY")
	tcg_client := tcg.NewClient(API_KEY)

	writer, err := createCsvWriter("./data/hashes.csv")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer writer.Flush()

	numWorkers := 4
	numPages := 10
	jobs := make(chan int, numPages)
	results := make(chan [][]string, numPages)

	// Start workers
	for w := 0; w < numWorkers; w++ {
		go func() {
			for page := range jobs {
				hashes := hashCards(page, tcg_client)
				results <- hashes
			}
		}()
	}

	// Send jobs
	for page := 1; page <= numPages; page++ {
		jobs <- page
	}
	close(jobs)

	// Collect results
	var hashes [][]string
	for i := 0; i < numPages; i++ {
		pageHashes := <-results
		hashes = append(hashes, pageHashes...)
	}

	for _, hash := range hashes {
		writeRecord(writer, hash)
	}

	elapsed := time.Since(start)
	fmt.Printf("Processing took %s\n", elapsed)
}
