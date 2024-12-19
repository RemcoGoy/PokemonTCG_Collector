package supabase

import (
	"os"

	"github.com/nedpals/supabase-go"
)

func NewSupabaseClient() *supabase.Client {
	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_KEY")

	supabase := supabase.CreateClient(url, key)
	return supabase
}
