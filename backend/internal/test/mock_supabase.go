package test

import (
	"github.com/supabase-community/gotrue-go"
	"github.com/supabase-community/gotrue-go/types"
	"github.com/supabase-community/supabase-go"
)

type MockSupabaseFactory struct {
	Auth gotrue.Client
}

func NewMockSupabaseFactory(auth gotrue.Client) *MockSupabaseFactory {
	return &MockSupabaseFactory{Auth: auth}
}

func (f *MockSupabaseFactory) CreateClient() *supabase.Client {
	if f.Auth == nil {
		return NewMockSupabaseClient(&MockAuth{})
	} else {
		return NewMockSupabaseClient(f.Auth)
	}
}

func (f *MockSupabaseFactory) CreateAuthenticatedClient(token string) *supabase.Client {
	client := f.CreateClient()
	session := types.Session{}
	session.AccessToken = token
	client.UpdateAuthSession(session)
	return client
}

func (f *MockSupabaseFactory) CreateAdminClient() *supabase.Client {
	return f.CreateClient()
}

type MockSupabaseClient struct {
	*supabase.Client
}

func NewMockSupabaseClient(auth gotrue.Client) *supabase.Client {
	return &supabase.Client{
		Auth: auth,
	}
}
