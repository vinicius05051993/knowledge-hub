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

	router.Handle(
		"POST /documents/upsert",
		middleware.AuthMiddleware(
			s.apiKeyService,
			http.HandlerFunc(
				s.documentHandler.UpsertDocuments,
			),
		),
	)

	router.Handle(
		"DELETE /documents",
		middleware.AuthMiddleware(
			s.apiKeyService,
			http.HandlerFunc(
				s.documentHandler.DeleteDocuments,
			),
		),
	)

	router.Handle(
		"POST /search",
		middleware.AuthMiddleware(
			s.apiKeyService,
			http.HandlerFunc(
				s.documentHandler.Search,
			),
		),
	)

	return router
}