package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/supabase-community/gotrue-go/types"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]any)
	var loginRequest LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		resp["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
	} else {
		token, err := s.SupabaseFactory.CreateClient().Auth.SignInWithEmailPassword(loginRequest.Email, loginRequest.Password)
		if err != nil {
			resp["error"] = err.Error()
			w.WriteHeader(http.StatusBadRequest)
		} else {
			resp["token"] = token.AccessToken
			w.Header().Set("Content-Type", "application/json")
		}
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		resp["error"] = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]any)
	if r.Header.Get("Authorization") == "" {
		resp["error"] = "Missing Authorization header"
		w.WriteHeader(http.StatusUnauthorized)
		jsonResp, _ := json.Marshal(resp)
		w.Write(jsonResp)
		return
	}

	token := r.Header.Get("Authorization")
	token = token[7:] // Remove "Bearer " prefix
	err := s.SupabaseFactory.CreateAuthenticatedClient(token).Auth.Logout()

	if err != nil {
		resp["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		resp["error"] = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) Signup(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]any)
	var registerRequest RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&registerRequest)
	if err != nil {
		resp["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
	} else {
		user, err := s.SupabaseFactory.CreateClient().Auth.Signup(types.SignupRequest{
			Email:    registerRequest.Email,
			Password: registerRequest.Password,
		})
		if err != nil {
			resp["error"] = err.Error()
			w.WriteHeader(http.StatusBadRequest)
		} else {
			resp["email"] = user.Email
			w.Header().Set("Content-Type", "application/json")
		}
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		resp["error"] = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]any)
	if r.Header.Get("Authorization") == "" {
		resp["error"] = "Missing Authorization header"
		w.WriteHeader(http.StatusUnauthorized)
		jsonResp, _ := json.Marshal(resp)
		w.Write(jsonResp)
		return
	}

	token := r.Header.Get("Authorization")
	token = token[7:] // Remove "Bearer " prefix
	user, err := s.SupabaseFactory.CreateAuthenticatedClient(token).Auth.GetUser()

	if err != nil {
		resp["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
	} else {
		resp["id"] = user.ID
		resp["email"] = user.Email
		resp["created_at"] = user.CreatedAt
		resp["role"] = user.Role
		w.Header().Set("Content-Type", "application/json")
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		resp["error"] = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	_, _ = w.Write(jsonResp)
}

func AuthRouter(s *Server) chi.Router {
	r := chi.NewRouter()
	r.Post("/login", s.Login)
	r.Post("/signup", s.Signup)
	r.Post("/logout", s.Logout)
	r.Get("/user", s.GetUser)
	return r
}
