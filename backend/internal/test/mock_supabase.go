package test

import (
	"github.com/nedpals/supabase-go"
)

// MockAuth implements the necessary auth methods
type MockAuth struct {
	*supabase.Auth
}

func NewMockSupabaseClient() *supabase.Client {
	return &supabase.Client{
		Auth: &MockAuth{},
	}
}
