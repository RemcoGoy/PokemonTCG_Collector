package supabase

import (
	"log"
	"os"

	"github.com/supabase-community/gotrue-go/types"
	"github.com/supabase-community/supabase-go"
)

type SupabaseFactory struct {
	url       string
	key       string
	admin_key string
}

type SupabaseFactoryInterface interface {
	CreateClient() *supabase.Client
	CreateAuthenticatedClient(token string) *supabase.Client
	CreateAdminClient() *supabase.Client
}

var _ SupabaseFactoryInterface = (*SupabaseFactory)(nil)

func NewSupabaseFactory() *SupabaseFactory {
	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_KEY")
	admin_key := os.Getenv("SUPABASE_ADMIN_KEY")

	return &SupabaseFactory{
		url:       url,
		key:       key,
		admin_key: admin_key,
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

func (f *SupabaseFactory) CreateAdminClient() *supabase.Client {
	client, err := supabase.NewClient(f.url, f.admin_key, &supabase.ClientOptions{})
	if err != nil {
		log.Fatalf("Error creating Supabase client: %v", err)
	}
	session := types.Session{}
	session.AccessToken = f.admin_key
	client.UpdateAuthSession(session)
	return client
}
