package integration

import (
	"testing"

	"indexer/internal/config"
	"indexer/internal/database"

	"github.com/jmoiron/sqlx"
)

func createTestConfig() *config.Config {

	return &config.Config{
		MySQLHost:     "localhost",
		MySQLPort:     "3306",
		MySQLDatabase: "indexer",
		MySQLUser:     "root",
		MySQLPassword: "root",

		OpenSearchHost: "localhost",
		OpenSearchPort: "9200",
	}
}

func createDB(
	t *testing.T,
) *sqlx.DB {

	cfg := createTestConfig()

	db, err := database.NewMySQL(cfg)

	if err != nil {
		t.Fatal(err)
	}

	return db
}

func cleanupTestData(
	t *testing.T,
	db *sqlx.DB,
) {

	t.Cleanup(func() {

		db.Exec(
			"DELETE FROM api_keys WHERE namespace = 'test'",
		)

		db.Exec(
			"DELETE FROM documents WHERE namespace = 'test'",
		)
	})
}