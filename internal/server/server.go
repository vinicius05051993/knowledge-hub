package server

import (
	"indexer/internal/apikeys"

	"github.com/jmoiron/sqlx"
)

type Server struct {
	db *sqlx.DB

	apiKeyService *apikeys.Service
}