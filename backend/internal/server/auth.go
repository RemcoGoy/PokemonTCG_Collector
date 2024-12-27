package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/supabase-community/gotrue-go/types"

	t "backend/internal/types"
	"backend/internal/utils"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]any)
	var loginRequest LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := s.SupabaseFactory.CreateClient().Auth.SignInWithEmailPassword(loginRequest.Email, loginRequest.Password)
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp["token"] = token.AccessToken
	w.Header().Set("Content-Type", "application/json")

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]any)
	if r.Header.Get("Authorization") == "" {
		utils.JSONError(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	token := r.Header.Get("Authorization")
	token = token[7:] // Remove "Bearer " prefix
	err := s.SupabaseFactory.CreateAuthenticatedClient(token).Auth.Logout()

	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) Signup(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]any)

	var registerRequest RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&registerRequest)
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.SupabaseFactory.CreateClient().Auth.Signup(types.SignupRequest{
		Email:    registerRequest.Email,
		Password: registerRequest.Password,
	})
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	profile := t.Profile{
		ID:       user.ID,
		Username: registerRequest.Username,
	}
	_, _, err = s.SupabaseFactory.CreateAdminClient().From("profile").Insert(profile, true, "", "", "exact").Execute()
	if err != nil {
		del_err := s.SupabaseFactory.CreateAdminClient().Auth.AdminDeleteUser(types.AdminDeleteUserRequest{UserID: user.ID})
		if del_err != nil {
			utils.JSONError(w, del_err.Error(), http.StatusBadRequest)
			return
		}

		utils.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp["email"] = user.Email
	w.Header().Set("Content-Type", "application/json")

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jsonResp)
}

func AuthRouter(s *Server) chi.Router {
	r := chi.NewRouter()
	r.Post("/login", s.Login)
	r.Post("/signup", s.Signup)
	r.Post("/logout", s.Logout)
	return r
}
