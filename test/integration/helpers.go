package integration

import (
	"context"
	"testing"

	"indexer/internal/apikeys"
	"indexer/internal/config"
	"indexer/internal/database"
	"indexer/internal/documentfilters"
	"indexer/internal/documents"
	"indexer/internal/opensearch"

	"github.com/jmoiron/sqlx"
)

func createTestConfig() *config.Config {

	return &config.Config{
		DBHost:     "localhost",
		DBPort:     "1433",
		DBName:     "indexer",
		DBUser:     "sa",
		DBPassword: "StrongPassword123!",
		OpenSearchUrl: "http://localhost:9200",
	}
}

func createDB(
	t *testing.T,
) *sqlx.DB {

	cfg := createTestConfig()

	db, err := database.New(*cfg)

	if err != nil {
		t.Fatal(err)
	}

	return db
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
			"Integration Test",
			"test",
		)

	if err != nil {
		t.Fatal(err)
	}

	return apiKey
}

func syncDocuments(
	t *testing.T,
	db *sqlx.DB,
) {

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

	syncService :=
		documents.NewSyncService(
			documentRepository,
			filterRepository,
			searchService,
		)

	err := syncService.ProcessPendingUpserts(
		context.Background(),
		1000,
	)

	if err != nil {
		t.Fatal(err)
	}

	err = syncService.ProcessPendingDeindexes(
		context.Background(),
		1000,
	)

	if err != nil {
		t.Fatal(err)
	}

	err = syncService.ProcessPendingDeletes(
		context.Background(),
		1000,
	)

	if err != nil {
		t.Fatal(err)
	}
}