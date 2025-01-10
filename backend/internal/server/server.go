package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"backend/internal/db"
	sb "backend/internal/supabase"

	tcg "github.com/PokemonTCG/pokemon-tcg-sdk-go-v2/pkg"
	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	SupabaseFactory sb.SupabaseFactoryInterface
	DbConnector     db.DbConnectorInterface
	TcgClient       tcg.Client
	port            int
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	supabaseFactory := sb.NewSupabaseFactory()

	// Download card_hashes.gob from Supabase storage
	client := supabaseFactory.CreateAdminClient()
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

	NewServer := &Server{
		SupabaseFactory: supabaseFactory,
		DbConnector:     db.NewDbConnector(supabaseFactory),
		TcgClient:       tcg.NewClient(os.Getenv("TCG_API_KEY")),
		port:            port,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
