package server

import (
	"indexer/internal/apikeys"
	"indexer/internal/documents"

	"github.com/jmoiron/sqlx"
)

type Server struct {
	db *sqlx.DB

	apiKeyService *apikeys.Service

	documentHandler *documents.Handler
}