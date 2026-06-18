package database

import (
	"fmt"
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/microsoft/go-mssqldb"

	"indexer/internal/config"
)

func New(cfg config.Config) (*sqlx.DB, error) {

	query := url.Values{}
	query.Add("database", cfg.DBName)
	query.Add("encrypt", "disable")

	dsn := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(cfg.DBUser, cfg.DBPassword),
		Host:     fmt.Sprintf("%s:%s", cfg.DBHost, cfg.DBPort),
		RawQuery: query.Encode(),
	}

	db, err := sqlx.Connect("sqlserver", dsn.String())
	if err != nil {
		return nil, fmt.Errorf("connect sql server: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping sql server: %w", err)
	}

	return db, nil
}