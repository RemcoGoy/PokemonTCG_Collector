package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	sb "backend/internal/supabase"

	_ "github.com/joho/godotenv/autoload"
	"github.com/nedpals/supabase-go"
)

type Server struct {
	SupabaseClient *supabase.Client
	port           int
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		SupabaseClient: sb.NewSupabaseClient(),
		port:           port,
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
