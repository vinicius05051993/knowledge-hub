package database

import (
	"fmt"

	"indexer/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

func NewMySQL(cfg *config.Config) (*sqlx.DB, error) {

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.MySQLUser,
		cfg.MySQLPassword,
		cfg.MySQLHost,
		cfg.MySQLPort,
		cfg.MySQLDatabase,
	)

	return sqlx.Connect("mysql", dsn)
}