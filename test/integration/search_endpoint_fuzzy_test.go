package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"indexer/internal/apikeys"
	"indexer/internal/documents"
	"indexer/internal/opensearch"
	"indexer/internal/server"
	"indexer/internal/documentfilters"
)

func TestSearchEndpointFuzzy(
	t *testing.T,
) {

	ensureDocumentsIndex(t)

	db := createDB(t)

	defer db.Close()

	cleanupTestData(
		t,
		db,
	)

	cfg := createTestConfig()

	searchClient :=
		opensearch.NewClient(cfg)

	t.Cleanup(func() {

		_ = opensearch.DeleteDocument(
			context.Background(),
			searchClient,
			"test",
			"fuzzy-test",
		)
	})

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

	documentHandler :=
		documents.NewHandler(
			documentService,
		)

	apiKeyRepository :=
		apikeys.NewRepository(db)

	apiKeyService :=
		apikeys.NewService(
			apiKeyRepository,
		)

	apiKey, err :=
		apiKeyService.Create(
			t.Context(),
			"test",
			"Fuzzy Test",
		)

	if err != nil {
		t.Fatal(err)
	}

	err = documentService.Upsert(
		t.Context(),
		&documents.Document{
			Namespace:  "test",
			ExternalID: "fuzzy-test",
			Title:      "Vinicius Henrique",
			Text:       "Vinicius Henrique platform",
			Payload: []byte(`{
				"sku":"fuzzy-test"
			}`),
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	app :=
		server.NewApp(
			db,
			apiKeyService,
			documentHandler,
		)

	router := app.Router()

	req := httptest.NewRequest(
		http.MethodPost,
		"/search",
		strings.NewReader(`{
			"query":"vinixcius"
		}`),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	req.Header.Set(
		"X-API-Key",
		apiKey,
	)

	recorder :=
		httptest.NewRecorder()

	router.ServeHTTP(
		recorder,
		req,
	)

	if recorder.Code !=
		http.StatusOK {

		t.Fatalf(
			"expected 200 got %d",
			recorder.Code,
		)
	}

	var response []documents.SearchResponse

	err = json.Unmarshal(
		recorder.Body.Bytes(),
		&response,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(response) != 1 {

		t.Fatalf(
			"expected 1 result got %d",
			len(response),
		)
	}

	if response[0].ExternalID !=
		"fuzzy-test" {

		t.Fatalf(
			"unexpected document returned: %s",
			response[0].ExternalID,
		)
	}
}