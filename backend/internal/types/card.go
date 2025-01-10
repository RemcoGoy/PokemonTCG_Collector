package types

import (
	"time"

	"github.com/google/uuid"
)

type Card struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	CollectionID uuid.UUID `json:"collection_id"`
	TCGID        string    `json:"tcg_id"`
	UserID       uuid.UUID `json:"user_id"`
}

type CardHash struct {
	ID         string `json:"id"`
	Perceptual string `json:"perceptual"`
	Difference string `json:"difference"`
	Wavelet    string `json:"wavelet"`
	Color      string `json:"color"`
	TCGID      string `json:"tcg_id"`
}

type CardHashGob struct {
	ID   string `json:"id"`
	Hash string `json:"hash"`
}
