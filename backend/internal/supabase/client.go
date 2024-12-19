package supabase

import (
	"os"

	supa "github.com/nedpals/supabase-go"
)

func NewSupabaseClient() *supa.Client {
	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_KEY")

	supabase := supa.CreateClient(url, key)
	return supabase
}
