package server

import (
	"indexer/internal/apikeys"
	"indexer/internal/documents"

	"github.com/jmoiron/sqlx"
)

func NewApp(
	db *sqlx.DB,
	apiKeyService *apikeys.Service,
	documentHandler *documents.Handler,
) *Server {

	return &Server{
		db: db,

		apiKeyService: apiKeyService,

		documentHandler: documentHandler,
	}
}