package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"backend/internal/db"
	sb "backend/internal/supabase"
	"backend/internal/types"
	"backend/internal/utils"

	tcg "github.com/PokemonTCG/pokemon-tcg-sdk-go-v2/pkg"
	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	SupabaseFactory sb.SupabaseFactoryInterface
	DbConnector     db.DbConnectorInterface
	TcgClient       tcg.Client
	port            int
	CARD_HASHES     []types.CardHashGob
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	supabaseFactory := sb.NewSupabaseFactory()

	utils.DownloadHashes(supabaseFactory)

	NewServer := &Server{
		SupabaseFactory: supabaseFactory,
		DbConnector:     db.NewDbConnector(supabaseFactory),
		TcgClient:       tcg.NewClient(os.Getenv("TCG_API_KEY")),
		port:            port,
		CARD_HASHES:     utils.ReadCardHashesGob(),
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
