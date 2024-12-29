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

	user_client := s.SupabaseFactory.CreateAuthenticatedClient(token)

	data, count, err := user_client.From("collection").Select("*", "exact", false).Eq("user_id", userID).Execute()
	if err != nil {
		utils.JSONError(w, "Could not fetch collections", http.StatusBadRequest)
		return
	}

	var collections []types.Collection
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
	collection := r.Context().Value(types.CollectionData).(types.Collection)

	jsonResp, err := json.Marshal(collection)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) CreateCollection(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]any)

	var createCollectionRequest CreateCollectionRequest
	err := json.NewDecoder(r.Body).Decode(&createCollectionRequest)
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	token := r.Context().Value(types.JwtTokenKey).(string)
	userID := r.Context().Value(types.UserID).(string)

	collection := types.Collection{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UserID:    uuid.MustParse(userID),
		Name:      createCollectionRequest.Name,
	}
	err = s.DbConnector.CreateCollection(collection, token)
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp["name"] = createCollectionRequest.Name

	jsonResp, err := json.Marshal(resp)
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

		user_client := s.SupabaseFactory.CreateAuthenticatedClient(token)
		data, _, err := user_client.From("collection").Select("*", "exact", false).Eq("id", id).Eq("user_id", userID).Execute()
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		var collections []types.Collection
		err = json.Unmarshal(data, &collections)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), types.CollectionID, id)
		ctx = context.WithValue(ctx, types.CollectionData, collections[0])
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
