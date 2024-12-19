package test

import (
	"github.com/nedpals/supabase-go"
)

func NewMockSupabaseClient() *supabase.Client {
	return &supabase.Client{
		Auth: &supabase.Auth{},
	}
}
