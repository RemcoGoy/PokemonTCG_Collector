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

// ListCollectionsHandler - Lists all collections for a user
//
//	@Summary		List all collections for a user
//	@Description	List all collections for a user
//	@Tags			Collection
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	types.ListCollectionsResponse
//	@Failure		400	{object}	types.ErrorResponse
//	@Router			/collection [get]
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

// GetCollectionHandler - Gets a collection for a user
//
//	@Summary		Get a collection for a user
//	@Description	Get a collection for a user by ID
//	@Tags			Collection
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Collection ID"
//	@Success		200	{object}	types.Collection
//	@Failure		400	{object}	types.ErrorResponse
//	@Router			/collection/{id} [get]
func (s *Server) GetCollection(w http.ResponseWriter, r *http.Request) {
	collection := r.Context().Value(types.CollectionData).(types.Collection)

	jsonResp, err := json.Marshal(collection)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jsonResp)
}

// CreateCollectionHandler - Creates a collection for a user
//
//	@Summary		Create a collection for a user
//	@Description	Create a collection for a user
//	@Tags			Collection
//	@Accept			json
//	@Produce		json
//	@Param			createCollectionRequest	body		types.CreateCollectionRequest	true	"Create collection request"
//	@Success		200	{object}	types.Collection
//	@Failure		400	{object}	types.ErrorResponse
//	@Router			/collection [post]
func (s *Server) CreateCollection(w http.ResponseWriter, r *http.Request) {
	var createCollectionRequest types.CreateCollectionRequest
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
		utils.JSONError(w, "error creating collection", http.StatusBadRequest)
		return
	}

	jsonResp, err := json.Marshal(collection)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jsonResp)
}

// UpdateCollectionHandler - Updates a collection for a user
//
//	@Summary		Update a collection for a user
//	@Description	Update a collection for a user by ID
//	@Tags			Collection
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Collection ID"
//	@Param			updateCollectionRequest	body		types.UpdateCollectionRequest	true	"Update collection request"
//	@Success		200	{object}	types.Collection
//	@Failure		400	{object}	types.ErrorResponse
//	@Router			/collection/{id} [put]
func (s *Server) UpdateCollection(w http.ResponseWriter, r *http.Request) {

}

// DeleteCollectionHandler - Deletes a collection for a user
//
//	@Summary		Delete a collection for a user
//	@Description	Delete a collection for a user by ID
//	@Tags			Collection
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Collection ID"
//	@Success		200	{object}	types.Collection
//	@Failure		400	{object}	types.ErrorResponse
//	@Router			/collection/{id} [delete]
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
