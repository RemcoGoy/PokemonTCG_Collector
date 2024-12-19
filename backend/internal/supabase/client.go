package supabase

import (
	"log"
	"os"

	"github.com/supabase-community/supabase-go"
)

type SupabaseClient struct {
	Client *supabase.Client
}

func NewSupabaseClient() *SupabaseClient {
	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_KEY")

	client, err := supabase.NewClient(url, key, &supabase.ClientOptions{})
	if err != nil {
		log.Fatalf("Failed to create Supabase client: %v", err)
	}

	return &SupabaseClient{Client: client}
}
