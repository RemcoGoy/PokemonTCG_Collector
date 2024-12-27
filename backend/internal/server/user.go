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

	data, count, err := user_client.From("profile").Select("*", "exact", false).Eq("id", user.ID.String()).Execute()
	if err != nil || count != 1 {
		utils.JSONError(w, "Could not fetch user profile", http.StatusBadRequest)
		return
	}

	var profiles []t.Profile
	err = json.Unmarshal(data, &profiles)
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	profile := profiles[0]
	resp["profile"] = profile
	resp["id"] = user.ID
	resp["email"] = user.Email
	resp["created_at"] = user.CreatedAt
	resp["role"] = user.Role
	w.Header().Set("Content-Type", "application/json")

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
