package server

import (
	"backend/internal/middleware"
	t "backend/internal/types"
	"backend/internal/utils"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) ListCollections(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]any)
	token := r.Context().Value(middleware.JwtTokenKey).(string)
	userID := r.Context().Value(middleware.UserID).(string)

	user_client := s.SupabaseFactory.CreateAuthenticatedClient(token)

	data, count, err := user_client.From("collection").Select("*", "exact", false).Eq("user_id", userID).Execute()
	if err != nil {
		utils.JSONError(w, "Could not fetch collections", http.StatusBadRequest)
		return
	}

	var collections []t.Collection
	err = json.Unmarshal(data, &collections)
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp["data"] = collections
	resp["count"] = count

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) GetCollection(w http.ResponseWriter, r *http.Request) {

}

func CollectionRouter(s *Server) chi.Router {
	r := chi.NewRouter()
	r.Get("/", s.ListCollections)
	r.Get("/{id}", s.GetCollection)
	return r
}
