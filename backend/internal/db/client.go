package db

import (
	"backend/internal/supabase"
	"backend/internal/types"
)

type DbConnector struct {
	supabaseFactory supabase.SupabaseFactoryInterface
}

type DbConnectorInterface interface {
	CreateCollection(collection types.Collection, token string) error
}

var _ DbConnectorInterface = (*DbConnector)(nil)

func NewDbConnector(supabaseFactory supabase.SupabaseFactoryInterface) *DbConnector {
	return &DbConnector{
		supabaseFactory: supabaseFactory,
	}
}
