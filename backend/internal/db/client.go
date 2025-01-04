package db

import (
	"backend/internal/supabase"
	"backend/internal/types"
)

type DbConnector struct {
	supabaseFactory supabase.SupabaseFactoryInterface
}

type DbConnectorInterface interface {
	// Profile
	CreateProfile(profile types.Profile) error
	GetProfile(id string, token string) (types.Profile, error)

	// Collection
	AdminCreateCollection(collection types.Collection) error
	CreateCollection(collection types.Collection, token string) (types.Collection, error)
	GetCollection(id string, userID string, token string) (types.Collection, error)
	ListCollections(userID string, token string) ([]types.Collection, int64, error)
	UpdateCollection(id string, update types.UpdateCollectionRequest, token string) (types.Collection, error)
	DeleteCollection(id string, userID string, token string) error

	// Card
	CreateCard(card types.Card, token string) (types.Card, error)
	GetCard(id string, userID string, token string) (types.Card, error)
	ListCards(userID string, token string) ([]types.Card, int64, error)
}

var _ DbConnectorInterface = (*DbConnector)(nil)

func NewDbConnector(supabaseFactory supabase.SupabaseFactoryInterface) *DbConnector {
	return &DbConnector{
		supabaseFactory: supabaseFactory,
	}
}
