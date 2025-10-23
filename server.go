package main

import (
	bo_v1 "lexicon/bo-api/beneficiary_ownership/v1"
	middlewares "lexicon/bo-api/middlewares"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"
)

type LexiconBOServer struct {
	router *chi.Mux
	cfg    config
}

func NewLexiconBOServer(cfg config) (*LexiconBOServer, error) {

	r := chi.NewRouter()

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://bo.lexicon.id", "https://beneficialowner.lexicon.id", "http://localhost:3000"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-API-KEY", "X-ACCESS-TIME", "X-REQUEST-SIGNATURE", "X-API-USER", "X-REQUEST-IDENTITY"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// // Set a timeout value on the request context (ctx), that will signal
	// // through ctx.Done() that the request has timed out and further
	// // processing should be stopped.
	r.Use(middleware.Timeout(2 * time.Minute))

	server := &LexiconBOServer{
		router: r,
		cfg:    cfg,
	}
	return server, nil

}

func (s *LexiconBOServer) setupRoute() {
	r := s.router
	cfg := s.cfg
	r.Route("/v1", func(r chi.Router) {
		r.Use(middlewares.AccessTime())
		r.Use(middlewares.ApiKey(cfg.BackendApiKey, cfg.ServerSalt))
		r.Use(middlewares.RequestSignature(cfg.ServerSalt))
		r.Mount("/beneficiary-ownership", bo_v1.Router())
	})
}

func (s *LexiconBOServer) start() {
	r := s.router
	cfg := s.cfg
	log.Info().Msg("Starting up server...")

	if err := http.ListenAndServe(cfg.Listen.Addr(), r); err != nil {
		log.Fatal().Err(err).Msg("Failed to start the server")
		return
	}

	log.Info().Msg("Server Stopped")
}
