package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"

	sb "backend/internal/supabase"
)

func DownloadHashes(factory sb.SupabaseFactoryInterface) {
	// Download card_hashes.gob from Supabase storage
	client := factory.CreateAdminClient()
	url, err := client.Storage.CreateSignedUrl("hashes", "card_hashes.gob", 60)
	if err != nil {
		fmt.Printf("Error creating signed url: %v\n", err)
	}
	resp, err := http.Get(url.SignedURL)
	if err != nil {
		fmt.Printf("Error downloading card hashes: %v\n", err)
	}
	defer resp.Body.Close()

	// Create local file
	out, err := os.OpenFile("./data/card_hashes.gob", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Printf("Error creating local file: %v\n", err)
	}
	defer out.Close()

	// Copy downloaded content to local file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("Error copying content to file: %v\n", err)
	}
}
