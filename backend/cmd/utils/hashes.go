package main

import (
	"backend/internal/utils"
	"encoding/csv"
	"fmt"
	"image"
	"net/http"
	"os"

	tcg "github.com/PokemonTCG/pokemon-tcg-sdk-go-v2/pkg"
	"github.com/PokemonTCG/pokemon-tcg-sdk-go-v2/pkg/request"
	"github.com/joho/godotenv"
	"github.com/schollz/progressbar/v3"
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

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	API_KEY := os.Getenv("TCG_API_KEY")
	tcg_client := tcg.NewClient(API_KEY)

	cards, err := tcg_client.GetCards(request.PageSize(10), request.Page(1))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	bar := progressbar.Default(int64(len(cards)))

	writer, err := createCsvWriter("./data/hashes.csv")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer writer.Flush()

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

		writeRecord(writer, []string{card.ID, hash.ToString()})
		bar.Add(1)
	}
}
