package integration

import (
	"context"
	"testing"

	"indexer/internal/apikeys"
	"indexer/internal/config"
	"indexer/internal/database"
	"indexer/internal/documents"
	"indexer/internal/opensearch"
	"indexer/internal/documentfilters"

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

func createAPIKeyService(
	t *testing.T,
	db *sqlx.DB,
) *apikeys.Service {

	repository :=
		apikeys.NewRepository(db)

	return apikeys.NewService(
		repository,
	)
}

func createDocumentHandler(
	t *testing.T,
	db *sqlx.DB,
) *documents.Handler {

	cfg := createTestConfig()

	searchClient :=
		opensearch.NewClient(cfg)

	searchService :=
		opensearch.NewService(
			searchClient,
		)

	documentRepository :=
		documents.NewRepository(db)

	filterRepository :=
		documentfilters.NewRepository(
			db,
		)

	documentService :=
		documents.NewService(
			documentRepository,
			filterRepository,
			searchService,
		)

	return documents.NewHandler(
		documentService,
	)
}

func createTestAPIKey(
	t *testing.T,
	db *sqlx.DB,
) string {

	repository :=
		apikeys.NewRepository(db)

	service :=
		apikeys.NewService(
			repository,
		)

	apiKey, err :=
		service.Create(
			context.Background(),
			"test",
			"Integration Test",
		)

	if err != nil {
		t.Fatal(err)
	}

	return apiKey
}