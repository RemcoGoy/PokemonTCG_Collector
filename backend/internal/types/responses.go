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
