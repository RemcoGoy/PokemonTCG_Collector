package server

import (
	"backend/internal/utils"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ScanHandler - Scans a card and returns the TCGID
//
//	@Summary		Scans a card and returns the TCGID
//	@Description	Scans a card and returns the TCGID
//	@Tags			Scan
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			card	formData	file	true	"Card image"
//	@Success		200		{object}	types.ScanResponse
//	@Failure		400		{object}	types.ErrorResponse
//	@Router			/scan [post]
func (s *Server) Scan(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(200 << 20); err != nil { // 200 MB max memory
		utils.JSONError(w, "File too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("card")
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	hash, err := utils.PhashImage(file, header)
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cardHash, err := utils.FindClosestCard(hash)
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	card, err := utils.GetCardData(cardHash.TCGID, s.TcgClient)
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(card)
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jsonResp)
}

func ScanRouter(s *Server) chi.Router {
	r := chi.NewRouter()
	r.Post("/", s.Scan)
	return r
}
