package db

import (
	"backend/internal/types"
	"encoding/json"
	"errors"
)

func (d *DbConnector) AdminCreateCollection(collection types.Collection) error {
	sb_client := d.supabaseFactory.CreateAdminClient()
	_, _, err := sb_client.From("collection").Insert(collection, false, "", "", "exact").Execute()
	if err != nil {
		return err
	}

	return nil
}

func (d *DbConnector) CreateCollection(collection types.Collection, token string) (types.Collection, error) {
	sb_client := d.supabaseFactory.CreateAuthenticatedClient(token)
	data, _, err := sb_client.From("collection").Insert(collection, false, "", "", "exact").Execute()
	if err != nil {
		return types.Collection{}, err
	}

	var c []types.Collection
	err = json.Unmarshal(data, &c)
	if err != nil {
		return types.Collection{}, err
	}

	return c[0], nil
}

func (d *DbConnector) GetCollection(id string, userID string, token string) (types.CollectionWithCards, error) {
	sb_client := d.supabaseFactory.CreateAuthenticatedClient(token)
	data, count, err := sb_client.From("collection").Select("*, card(*)", "exact", false).Eq("user_id", userID).Eq("id", id).Execute()
	if err != nil || count == 0 {
		return types.CollectionWithCards{}, errors.New("collection not found")
	}

	var collections []types.CollectionWithCards
	err = json.Unmarshal(data, &collections)
	if err != nil {
		return types.CollectionWithCards{}, err
	}

	return collections[0], nil
}

func (d *DbConnector) ListCollections(userID string, token string) ([]types.CollectionWithCards, int64, error) {
	sb_client := d.supabaseFactory.CreateAuthenticatedClient(token)

	data, count, err := sb_client.From("collection").Select("*, card(*)", "exact", false).Eq("user_id", userID).Execute()
	if err != nil {
		return nil, 0, err
	}

	var collections []types.CollectionWithCards
	err = json.Unmarshal(data, &collections)
	if err != nil {
		return nil, 0, err
	}

	return collections, count, nil
}

func (d *DbConnector) UpdateCollection(id string, update types.UpdateCollectionRequest, token string) (types.Collection, error) {
	sb_client := d.supabaseFactory.CreateAuthenticatedClient(token)
	data, _, err := sb_client.From("collection").Update(update, "", "exact").Eq("id", id).Execute()

	if err != nil {
		return types.Collection{}, err
	}

	var c []types.Collection
	err = json.Unmarshal(data, &c)
	if err != nil {
		return types.Collection{}, err
	}

	return c[0], nil
}

func (d *DbConnector) DeleteCollection(id string, userID string, token string) error {
	sb_client := d.supabaseFactory.CreateAuthenticatedClient(token)
	_, _, err := sb_client.From("collection").Delete("", "exact").Eq("id", id).Execute()
	if err != nil {
		return err
	}

	return nil
}
