package test

import (
	"backend/internal/db"
	"backend/internal/types"
)

type MockDbConnector struct {
}

var _ db.DbConnectorInterface = (*MockDbConnector)(nil)

func NewMockDbConnector() *MockDbConnector {
	return &MockDbConnector{}
}

func (m *MockDbConnector) CreateProfile(profile types.Profile) error {
	return nil
}

func (m *MockDbConnector) GetProfile(id string, token string) (types.Profile, error) {
	return types.Profile{}, nil
}

func (m *MockDbConnector) CreateCollection(collection types.Collection, token string) error {
	return nil
}

func (m *MockDbConnector) GetCollection(id string, userID string, token string) (types.Collection, error) {
	return types.Collection{}, nil
}

func (m *MockDbConnector) ListCollections(userID string, token string) ([]types.Collection, int64, error) {
	return []types.Collection{}, 0, nil
}
