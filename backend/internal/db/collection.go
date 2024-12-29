package db

import (
	"backend/internal/types"
)

func (d *DbConnector) CreateCollection(collection types.Collection, token string) error {
	sb_client := d.supabaseFactory.CreateAuthenticatedClient(token)
	_, _, err := sb_client.From("collection").Insert(collection, false, "", "", "exact").Execute()
	return err
}
