package db

import (
	"backend/internal/types"
	"encoding/json"
)

func (d *DbConnector) CreateCollection(collection types.Collection, token string) (types.Collection, error) {
	sb_client := d.supabaseFactory.CreateAuthenticatedClient(token)
	data, _, err := sb_client.From("collection").Insert(collection, false, "", "", "exact").Execute()
	if err != nil {
		return types.Collection{}, err
	}

	var c types.Collection
	err = json.Unmarshal(data, &c)
	if err != nil {
		return types.Collection{}, err
	}

	return c, nil
}

func (d *DbConnector) GetCollection(id string, userID string, token string) (types.Collection, error) {
	sb_client := d.supabaseFactory.CreateAuthenticatedClient(token)
	data, _, err := sb_client.From("collection").Select("*", "exact", false).Eq("user_id", userID).Execute()
	if err != nil {
		return types.Collection{}, err
	}

	var collections []types.Collection
	err = json.Unmarshal(data, &collections)
	if err != nil {
		return types.Collection{}, err
	}

	return collections[0], nil
}

func (d *DbConnector) ListCollections(userID string, token string) ([]types.Collection, int64, error) {
	sb_client := d.supabaseFactory.CreateAuthenticatedClient(token)

	data, count, err := sb_client.From("collection").Select("*", "exact", false).Eq("user_id", userID).Execute()
	if err != nil {
		return nil, 0, err
	}

	var collections []types.Collection
	err = json.Unmarshal(data, &collections)
	if err != nil {
		return nil, 0, err
	}

	return collections, count, nil
}
