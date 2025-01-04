package types

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type CreateCollectionRequest struct {
	Name string `json:"name"`
}

type UpdateCollectionRequest struct {
	Name string `json:"name"`
}

type CreateCardRequest struct {
	CollectionID string `json:"collection_id"`
	TCGID        string `json:"tcg_id"`
}

type UpdateCardRequest struct {
	CollectionID string `json:"collection_id"`
	TCGID        string `json:"tcg_id"`
}
