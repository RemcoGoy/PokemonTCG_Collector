package db

import (
	"backend/internal/types"
	"encoding/json"
	"errors"
)

func (d *DbConnector) CreateProfile(profile types.Profile) error {
	sb_client := d.supabaseFactory.CreateAdminClient()
	_, _, err := sb_client.From("profile").Insert(profile, true, "", "", "exact").Execute()
	if err != nil {
		return err
	}

	return nil
}

func (d *DbConnector) GetProfile(id string, token string) (types.Profile, error) {
	sb_client := d.supabaseFactory.CreateAuthenticatedClient(token)

	data, count, err := sb_client.From("profile").Select("*", "exact", false).Eq("id", id).Execute()
	if err != nil {
		return types.Profile{}, err
	}

	if count != 1 {
		return types.Profile{}, errors.New("profile not found")
	}

	var profiles []types.Profile
	err = json.Unmarshal(data, &profiles)
	if err != nil {
		return types.Profile{}, err
	}

	profile := profiles[0]
	return profile, nil
}
