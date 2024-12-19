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
		token, err := s.SupabaseClient.Auth.SignInWithEmailPassword(loginRequest.Email, loginRequest.Password)
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

func (s *Server) Signup(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]any)
	var registerRequest RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&registerRequest)
	if err != nil {
		resp["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
	} else {
		user, err := s.SupabaseClient.Auth.Signup(types.SignupRequest{
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

func AuthRouter(s *Server) chi.Router {
	r := chi.NewRouter()
	r.Post("/login", s.Login)
	r.Post("/signup", s.Signup)
	return r
}
