package database

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/jmoiron/sqlx"

	_ "github.com/microsoft/go-mssqldb"
	_ "github.com/microsoft/go-mssqldb/azuread"

	"indexer/internal/config"
)

func OpenDatabase(cfg config.Config) (*sql.DB, string, error) {

	var (
		driver string
		dsn    string
	)

	if cfg.UseManagedIdentity {

		driver = "azuresql"

		dsn = fmt.Sprintf(
			"sqlserver://%s:%s?database=%s&fedauth=ActiveDirectoryManagedIdentity&encrypt=true",
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBName,
		)

	} else {

		driver = "sqlserver"

		query := url.Values{}
		query.Add("database", cfg.DBName)
		query.Add("encrypt", "disable")

		u := &url.URL{
			Scheme: "sqlserver",
			User: url.UserPassword(
				cfg.DBUser,
				cfg.DBPassword,
			),
			Host: fmt.Sprintf(
				"%s:%s",
				cfg.DBHost,
				cfg.DBPort,
			),
			RawQuery: query.Encode(),
		}

		dsn = u.String()
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, "", fmt.Errorf("open sql server: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, "", fmt.Errorf("ping sql server: %w", err)
	}

	return db, driver, nil
}

func New(cfg config.Config) (*sqlx.DB, error) {

	db, _, err := OpenDatabase(cfg)

	if err != nil {
		return nil, err
	}

	return sqlx.NewDb(
		db,
		"sqlserver",
	), nil
}