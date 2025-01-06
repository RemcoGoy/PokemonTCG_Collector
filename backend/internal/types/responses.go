package types

type ErrorResponse struct {
	Error string `json:"error"`
}

type OkResponse struct {
	Ok bool `json:"ok"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type LogoutResponse struct {
}

type RegisterResponse struct {
	Email string `json:"email"`
}

type ListResponse struct {
	Data  []any `json:"data"`
	Count int   `json:"count"`
}

type ListCollectionsResponse struct {
	ListResponse
	Data []Collection `json:"data"`
}

type ListCardsResponse struct {
	ListResponse
	Data []Card `json:"data"`
}

type ScanResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Supertype   string   `json:"supertype"`
	Subtypes    []string `json:"subtypes"`
	Level       string   `json:"level"`
	HP          string   `json:"hp"`
	Types       []string `json:"types"`
	EvolvesFrom string   `json:"evolves_from"`
	EvolvesTo   []string `json:"evolves_to"`
}
