package db

import (
	"backend/internal/types"
	"encoding/json"
	"errors"
)

func (d *DbConnector) GetCard(id string, userID string, token string) (types.Card, error) {
	sb_client := d.supabaseFactory.CreateAuthenticatedClient(token)

	data, count, err := sb_client.From("card").Select("*", "exact", false).Eq("user_id", userID).Eq("id", id).Execute()
	if err != nil || count == 0 {
		return types.Card{}, errors.New("card not found")
	}

	var cards []types.Card
	err = json.Unmarshal(data, &cards)
	if err != nil {
		return types.Card{}, err
	}

	return cards[0], nil
}

func (d *DbConnector) ListCards(userID string, token string) ([]types.Card, int64, error) {
	sb_client := d.supabaseFactory.CreateAuthenticatedClient(token)

	data, count, err := sb_client.From("card").Select("*", "exact", false).Eq("user_id", userID).Execute()
	if err != nil {
		return []types.Card{}, 0, err
	}

	var cards []types.Card
	err = json.Unmarshal(data, &cards)
	if err != nil {
		return []types.Card{}, 0, err
	}

	return cards, count, nil
}

func (d *DbConnector) CreateCard(card types.Card, token string) (types.Card, error) {
	sb_client := d.supabaseFactory.CreateAuthenticatedClient(token)
	data, _, err := sb_client.From("card").Insert(card, false, "", "", "exact").Execute()
	if err != nil {
		return types.Card{}, err
	}

	var c []types.Card
	err = json.Unmarshal(data, &c)
	if err != nil {
		return types.Card{}, err
	}

	return c[0], nil
}

func (d *DbConnector) UpdateCard(id string, update types.UpdateCardRequest, token string) (types.Card, error) {
	sb_client := d.supabaseFactory.CreateAuthenticatedClient(token)
	data, _, err := sb_client.From("card").Update(update, "", "exact").Eq("id", id).Execute()
	if err != nil {
		return types.Card{}, err
	}

	var c []types.Card
	err = json.Unmarshal(data, &c)
	if err != nil {
		return types.Card{}, err
	}

	return c[0], nil
}

func (d *DbConnector) DeleteCard(id string, userID string, token string) error {
	return nil
}
