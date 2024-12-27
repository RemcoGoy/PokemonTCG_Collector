package types

import (
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
}
