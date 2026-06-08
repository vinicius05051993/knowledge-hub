package server

import (
	"indexer/internal/apikeys"

	"github.com/jmoiron/sqlx"
)

func NewApp(
	db *sqlx.DB,
	apiKeyService *apikeys.Service,
) *Server {

	return &Server{
		db:            db,
		apiKeyService: apiKeyService,
	}
}