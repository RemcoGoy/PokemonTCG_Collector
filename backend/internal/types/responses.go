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
