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

// GetCardHandler - Gets a card for a user
//
//	@Summary		Get a card for a user
//	@Description	Get a card for a user by ID
//	@Tags			Card
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Card ID"
//	@Success		200	{object}	types.Card
//	@Failure		400	{object}	types.ErrorResponse
//	@Router			/card/{id} [get]
func (s *Server) GetCard(w http.ResponseWriter, r *http.Request) {
	card := r.Context().Value(types.CardData).(types.Card)

	jsonResp, err := json.Marshal(card)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jsonResp)
}

// CreateCardHandler - Creates a card for a user
//
//	@Summary		Create a card for a user
//	@Description	Create a card for a user
//	@Tags			Card
//	@Accept			json
//	@Produce		json
//	@Param			body	body		types.CreateCardRequest	true	"Create Card Request"
//	@Success		200		{object}	types.Card
//	@Failure		400		{object}	types.ErrorResponse
//	@Router			/card [post]
func (s *Server) CreateCard(w http.ResponseWriter, r *http.Request) {
	var createCardRequest types.CreateCardRequest
	err := json.NewDecoder(r.Body).Decode(&createCardRequest)
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	token := r.Context().Value(types.JwtTokenKey).(string)
	userID := r.Context().Value(types.UserID).(string)

	card := types.Card{
		ID:           uuid.New(),
		CreatedAt:    time.Now(),
		CollectionID: uuid.MustParse(createCardRequest.CollectionID),
		TCGID:        createCardRequest.TCGID,
		UserID:       uuid.MustParse(userID),
	}
	card, err = s.DbConnector.CreateCard(card, token)
	if err != nil {
		utils.JSONError(w, "error creating card", http.StatusBadRequest)
		return
	}

	jsonResp, err := json.Marshal(card)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jsonResp)
}

// UpdateCardHandler - Updates a card for a user
//
//	@Summary		Update a card for a user
//	@Description	Update a card for a user by ID
//	@Tags			Card
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"Card ID"
//	@Param			body	body		types.UpdateCardRequest	true	"Update Card Request"
//	@Success		200		{object}	types.Card
//	@Failure		400		{object}	types.ErrorResponse
//	@Router			/card/{id} [put]
func (s *Server) UpdateCard(w http.ResponseWriter, r *http.Request) {
	prevCard := r.Context().Value(types.CardData).(types.Card)
	token := r.Context().Value(types.JwtTokenKey).(string)

	var updateCardRequest types.UpdateCardRequest
	err := json.NewDecoder(r.Body).Decode(&updateCardRequest)
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updateCardRequest.CollectionID == "" {
		updateCardRequest.CollectionID = prevCard.CollectionID.String()
	}

	newCard, err := s.DbConnector.UpdateCard(prevCard.ID.String(), updateCardRequest, token)
	if err != nil {
		utils.JSONError(w, "error updating card", http.StatusBadRequest)
		return
	}

	jsonResp, err := json.Marshal(newCard)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jsonResp)
}

// DeleteCardHandler - Deletes a card for a user
//
//	@Summary		Delete a card for a user
//	@Description	Delete a card for a user by ID
//	@Tags			Card
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Card ID"
//	@Success		200	{object}	types.Card
//	@Failure		400	{object}	types.ErrorResponse
//	@Router			/card/{id} [delete]
func (s *Server) DeleteCard(w http.ResponseWriter, r *http.Request) {
	card := r.Context().Value(types.CardData).(types.Card)
	token := r.Context().Value(types.JwtTokenKey).(string)

	err := s.DbConnector.DeleteCard(card.ID.String(), card.UserID.String(), token)
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResp, err := json.Marshal(card)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jsonResp)
}

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
