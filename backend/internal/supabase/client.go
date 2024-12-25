package supabase

import (
	"log"
	"os"

	"github.com/supabase-community/gotrue-go/types"
	"github.com/supabase-community/supabase-go"
)

type SupabaseFactory struct {
	url string
	key string
}

func NewSupabaseFactory() *SupabaseFactory {
	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_KEY")

	return &SupabaseFactory{
		url: url,
		key: key,
	}
}

func (f *SupabaseFactory) CreateClient() *supabase.Client {
	client, err := supabase.NewClient(f.url, f.key, &supabase.ClientOptions{})
	if err != nil {
		log.Fatalf("Error creating Supabase client: %v", err)
	}
	return client
}

func (f *SupabaseFactory) CreateAuthenticatedClient(token string) *supabase.Client {
	client := f.CreateClient()
	session := types.Session{}
	session.AccessToken = token
	client.UpdateAuthSession(session)
	return client
}
