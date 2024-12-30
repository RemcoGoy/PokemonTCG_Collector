package server

import (
	"backend/internal/types"
	"backend/internal/utils"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CreateCollectionRequest struct {
	Name string `json:"name"`
}

func (s *Server) ListCollections(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]any)
	token := r.Context().Value(types.JwtTokenKey).(string)
	userID := r.Context().Value(types.UserID).(string)

	collections, count, err := s.DbConnector.ListCollections(userID, token)
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusBadRequest)
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
	collection := r.Context().Value(types.CollectionData).(types.Collection)

	jsonResp, err := json.Marshal(collection)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) CreateCollection(w http.ResponseWriter, r *http.Request) {
	var createCollectionRequest CreateCollectionRequest
	err := json.NewDecoder(r.Body).Decode(&createCollectionRequest)
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	token := r.Context().Value(types.JwtTokenKey).(string)
	userID := r.Context().Value(types.UserID).(string)

	if createCollectionRequest.Name == "" {
		utils.JSONError(w, "name is required", http.StatusBadRequest)
		return
	}

	collection := types.Collection{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UserID:    uuid.MustParse(userID),
		Name:      createCollectionRequest.Name,
	}
	collection, err = s.DbConnector.CreateCollection(collection, token)
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResp, err := json.Marshal(collection)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) UpdateCollection(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) DeleteCollection(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) CollectionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		token := r.Context().Value(types.JwtTokenKey).(string)
		userID := r.Context().Value(types.UserID).(string)

		collection, err := s.DbConnector.GetCollection(id, userID, token)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), types.CollectionID, id)
		ctx = context.WithValue(ctx, types.CollectionData, collection)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CollectionRouter(s *Server) chi.Router {
	r := chi.NewRouter()
	r.Get("/", s.ListCollections)
	r.Post("/", s.CreateCollection)
	r.Route("/{id}", func(r chi.Router) {
		r.Use(s.CollectionCtx)
		r.Get("/", s.GetCollection)
		r.Delete("/", s.DeleteCollection)
		r.Put("/", s.UpdateCollection)
	})
	return r
}
