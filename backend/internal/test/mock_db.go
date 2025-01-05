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

func (m *MockDbConnector) AdminCreateCollection(collection types.Collection) error {
	return nil
}

func (m *MockDbConnector) CreateCollection(collection types.Collection, token string) (types.Collection, error) {
	return types.Collection{
		ID:   uuid.MustParse(COLLECTION_ID),
		Name: COLLECTION_NAME,
	}, nil
}

func (m *MockDbConnector) GetCollection(id string, userID string, token string) (types.CollectionWithCards, error) {
	return types.CollectionWithCards{
		Collection: types.Collection{
			ID:   uuid.MustParse(COLLECTION_ID),
			Name: COLLECTION_NAME,
		},
		Card: []types.Card{},
	}, nil
}

func (m *MockDbConnector) ListCollections(userID string, token string) ([]types.CollectionWithCards, int64, error) {
	return []types.CollectionWithCards{}, 0, nil
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

func (m *MockDbConnector) GetCard(id string, userID string, token string) (types.Card, error) {
	return types.Card{}, nil
}

func (m *MockDbConnector) ListCards(userID string, token string) ([]types.Card, int64, error) {
	return []types.Card{}, 0, nil
}

func (m *MockDbConnector) CreateCard(card types.Card, token string) (types.Card, error) {
	return types.Card{}, nil
}

func (m *MockDbConnector) UpdateCard(id string, update types.UpdateCardRequest, token string) (types.Card, error) {
	return types.Card{}, nil
}

func (m *MockDbConnector) DeleteCard(id string, userID string, token string) error {
	return nil
}

// Not Found Db Connector

type NotFoundDbConnector struct {
	MockDbConnector
}

func NewNotFoundDbConnector() *NotFoundDbConnector {
	return &NotFoundDbConnector{}
}

func (m *NotFoundDbConnector) GetCollection(id string, userID string, token string) (types.CollectionWithCards, error) {
	return types.CollectionWithCards{}, errors.New("collection not found")
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
