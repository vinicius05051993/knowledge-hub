package server

import (
	"net/http"

	"indexer/internal/middleware"
)

func (s *Server) Router() *http.ServeMux {

	router := http.NewServeMux()

	router.HandleFunc(
		"GET /health",
		HealthHandler,
	)

	router.Handle(
		"GET /ready",
		ReadyHandler(s.db),
	)

	router.Handle(
		"GET /protected",
		middleware.AuthMiddleware(
			s.apiKeyService,
			http.HandlerFunc(
				ProtectedHandler,
			),
		),
	)

	return router
}