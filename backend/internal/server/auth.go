package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/supabase-community/gotrue-go/types"

	t "backend/internal/types"
	"backend/internal/utils"
)

// LoginHandler - Logs in a user
//
//	@Summary		Login a user
//	@Description	Log a user in by email and password
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			loginRequest	body		types.LoginRequest	true	"Login request"
//	@Success		200				{object}	types.LoginResponse
//	@Failure		400				{object}	types.ErrorResponse
//	@Router			/auth/login [post]
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]any)
	var loginRequest t.LoginRequest
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

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jsonResp)
}

// LogoutHandler - Logs out a user
//
//	@Summary		Logout a user
//	@Description	Log a user out
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	types.LogoutResponse
//	@Failure		400	{object}	types.ErrorResponse
//	@Failure		401	{object}	types.ErrorResponse
//	@Router			/auth/logout [post]
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

// SignupHandler - Signs up a user
//
//	@Summary		Signup a user
//	@Description	Sign up a user by email, password and username
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			registerRequest	body		types.RegisterRequest	true	"Register request"
//	@Success		200				{object}	types.RegisterResponse
//	@Failure		400				{object}	types.ErrorResponse
//	@Router			/auth/signup [post]
func (s *Server) Signup(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]any)

	var registerRequest t.RegisterRequest
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
		ID:        user.ID,
		Username:  registerRequest.Username,
		CreatedAt: time.Now(),
	}

	creation_err := s.DbConnector.CreateProfile(profile)
	if creation_err == nil {
		creation_err = s.DbConnector.AdminCreateCollection(t.Collection{
			ID:        uuid.New(),
			Name:      "Default",
			UserID:    user.ID,
			CreatedAt: time.Now(),
		})
	}

	if creation_err != nil {
		del_err := s.SupabaseFactory.CreateAdminClient().Auth.AdminDeleteUser(types.AdminDeleteUserRequest{UserID: user.ID})
		if del_err != nil {
			utils.JSONError(w, del_err.Error(), http.StatusBadRequest)
			return
		}

		utils.JSONError(w, creation_err.Error(), http.StatusBadRequest)
		return
	}

	resp["email"] = user.Email

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
