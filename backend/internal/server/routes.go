package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "backend/docs"
	m "backend/internal/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Public routes
	r.Group(func(r chi.Router) {
		r.Mount("/swagger", httpSwagger.WrapHandler)
		r.Get("/ok", s.okHandler)
		r.Mount("/auth", AuthRouter(s))
	})

	// Private routes
	r.Group(func(r chi.Router) {
		r.Use(m.CheckJwtToken)
		r.Mount("/user", UserRouter(s))
		r.Mount("/collection", CollectionRouter(s))
		r.Mount("/card", CardRouter(s))
		r.Mount("/scan", ScanRouter(s))
	})

	return r
}

// OkHandler - Checks if the API is working
//
//	@Summary		This API can be used as health check for this application.
//	@Description	Tells if the API is working or not.
//	@Tags			Health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	types.OkResponse
//	@Router			/ok [get]
func (s *Server) okHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]bool)
	resp["ok"] = true

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
