package server

import (
	"backend/internal/types"
	"backend/internal/utils"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) Scan(w http.ResponseWriter, r *http.Request) {
	resp := types.ScanResponse{}

	if err := r.ParseMultipartForm(200 << 20); err != nil { // 200 MB max memory
		utils.JSONError(w, "File too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("card.png") // TODO: change this, but there's a bug in Scalar
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

	card, err := utils.FindClosestCard(hash)
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.Hash = hash.ToString()
	resp.CardHash = card

	jsonResp, err := json.Marshal(resp)
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
