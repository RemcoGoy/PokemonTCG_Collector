package test

import (
	"github.com/supabase-community/supabase-go"
)

func NewMockSupabaseClient() *supabase.Client {
	return &supabase.Client{}
}
