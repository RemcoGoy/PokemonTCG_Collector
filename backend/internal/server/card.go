package server

import (
	"backend/internal/types"
	"backend/internal/utils"
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ListCardsHandler - Lists all cards for a user
//
//	@Summary		List all cards for a user
//	@Description	List all cards for a user
//	@Tags			Card
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	types.ListCardsResponse
//	@Failure		400	{object}	types.ErrorResponse
//	@Router			/card [get]
func (s *Server) ListCards(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]any)
	token := r.Context().Value(types.JwtTokenKey).(string)
	userID := r.Context().Value(types.UserID).(string)

	cards, count, err := s.DbConnector.ListCards(userID, token)
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp["data"] = cards
	resp["count"] = count

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) GetCard(w http.ResponseWriter, r *http.Request) {}

func (s *Server) CreateCard(w http.ResponseWriter, r *http.Request) {}

func (s *Server) UpdateCard(w http.ResponseWriter, r *http.Request) {}

func (s *Server) DeleteCard(w http.ResponseWriter, r *http.Request) {}

func (s *Server) CardCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		token := r.Context().Value(types.JwtTokenKey).(string)
		userID := r.Context().Value(types.UserID).(string)

		card, err := s.DbConnector.GetCard(id, userID, token)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), types.CardID, id)
		ctx = context.WithValue(ctx, types.CardData, card)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CardRouter(s *Server) chi.Router {
	r := chi.NewRouter()
	r.Get("/", s.ListCards)
	r.Post("/", s.CreateCard)
	r.Route("/{id}", func(r chi.Router) {
		r.Use(s.CardCtx)
		r.Get("/", s.GetCard)
		r.Delete("/", s.DeleteCard)
		r.Put("/", s.UpdateCard)
	})
	return r
}
