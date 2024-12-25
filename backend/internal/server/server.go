package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	sb "backend/internal/supabase"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	SupabaseFactory sb.SupabaseFactoryInterface
	port            int
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		SupabaseFactory: sb.NewSupabaseFactory(),
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