package server

import (
	t "backend/internal/types"
	"backend/internal/utils"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]any)
	token := r.Context().Value(t.JwtTokenKey).(string)

	user_client := s.SupabaseFactory.CreateAuthenticatedClient(token)

	user, err := user_client.Auth.GetUser()
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	profile, err := s.DbConnector.GetProfile(user.ID.String(), token)
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp["profile"] = profile
	resp["id"] = user.ID
	resp["email"] = user.Email
	resp["created_at"] = user.CreatedAt
	resp["role"] = user.Role

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jsonResp)
}

func UserRouter(s *Server) chi.Router {
	r := chi.NewRouter()
	r.Get("/profile", s.GetUser)
	return r
}
