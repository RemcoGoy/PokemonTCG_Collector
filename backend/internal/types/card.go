package types

import (
	"time"

	"github.com/google/uuid"
)

type Card struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	CollectionID uuid.UUID `json:"collection_id"`
}
