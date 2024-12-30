package test

import (
	"backend/internal/db"
	"backend/internal/types"
	"errors"

	"github.com/google/uuid"
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

func (m *MockDbConnector) CreateCollection(collection types.Collection, token string) (types.Collection, error) {
	return types.Collection{
		ID:   uuid.MustParse(COLLECTION_ID),
		Name: COLLECTION_NAME,
	}, nil
}

func (m *MockDbConnector) GetCollection(id string, userID string, token string) (types.Collection, error) {
	return types.Collection{
		ID:   uuid.MustParse(COLLECTION_ID),
		Name: COLLECTION_NAME,
	}, nil
}

func (m *MockDbConnector) ListCollections(userID string, token string) ([]types.Collection, int64, error) {
	return []types.Collection{}, 0, nil
}

func (m *MockDbConnector) UpdateCollection(id string, update types.UpdateCollectionRequest, token string) (types.Collection, error) {
	return types.Collection{
		ID:   uuid.MustParse(COLLECTION_ID),
		Name: update.Name,
	}, nil
}

func (m *MockDbConnector) DeleteCollection(id string, userID string, token string) error {
	return nil
}

// Not Found Db Connector

type NotFoundDbConnector struct {
	MockDbConnector
}

func NewNotFoundDbConnector() *NotFoundDbConnector {
	return &NotFoundDbConnector{}
}

func (m *NotFoundDbConnector) GetCollection(id string, userID string, token string) (types.Collection, error) {
	return types.Collection{}, errors.New("collection not found")
}

// Duplicate Db Connector

type DuplicateDbConnector struct {
	MockDbConnector
}

func NewDuplicateDbConnector() *DuplicateDbConnector {
	return &DuplicateDbConnector{}
}

func (m *DuplicateDbConnector) CreateCollection(collection types.Collection, token string) (types.Collection, error) {
	return types.Collection{}, errors.New("(23505) duplicate key value violates unique constraint \"collection_name_key\"")
}
