package types

import (
	"time"

	"github.com/google/uuid"
)

type Collection struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	UserID    uuid.UUID `json:"user_id"`
}
