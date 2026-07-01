package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"indexer/internal/apikeys"
	"indexer/internal/documentfilters"
	"indexer/internal/documents"
	"indexer/internal/opensearch"
	"indexer/internal/server"
)

func TestSearchEndpointWithFilters(
	t *testing.T,
) {

	ensureDocumentsIndex(t)

	db := createDB(t)

	cfg := createTestConfig()

	searchClient :=
		opensearch.NewClient(cfg)

	apiKeyRepository :=
		apikeys.NewRepository(db)

	apiKeyService :=
		apikeys.NewService(
			apiKeyRepository,
		)

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

	t.Cleanup(func() {

		_ = opensearch.DeleteDocument(
			context.Background(),
			searchClient,
			"search-test",
			"filter-1",
		)

		_ = opensearch.DeleteDocument(
			context.Background(),
			searchClient,
			"search-test",
			"filter-2",
		)

		_ = filterRepository.DeleteByDocumentKeys(
			context.Background(),
			[]string{
				"search-test:filter-1",
				"search-test:filter-2",
			},
		)

		_ = documentService.DeleteByNamespace(
			context.Background(),
			"search-test",
		)

		_ = apiKeyService.DeleteByNamespace(
			context.Background(),
			"search-test",
		)

		_ = db.Close()
	})

	apiKey, err :=
		apiKeyService.Create(
			t.Context(),
			"Search Test",
			"search-test",
		)

	if err != nil {
		t.Fatal(err)
	}

	err = documentService.Upsert(
		t.Context(),
		&documents.Document{
			Namespace:  "search-test",
			ExternalID: "filter-1",
			Title:      "Magento",
			Text:       "Magento ecommerce",
			Payload:    []byte(`{"sku":"abc"}`),
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	err = documentService.Upsert(
		t.Context(),
		&documents.Document{
			Namespace:  "search-test",
			ExternalID: "filter-2",
			Title:      "Magento",
			Text:       "Magento ecommerce",
			Payload:    []byte(`{"sku":"xyz"}`),
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	syncDocuments(
		t,
		db,
	)

	documentHandler :=
		documents.NewHandler(
			documentService,
		)

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
			"query":"Magento",
			"filters":[
				{
					"field":"sku",
					"operator":"eq",
					"value":"abc",
					"value_type":"string"
				}
			]
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
		"filter-1" {

		t.Fatal(
			"wrong document returned",
		)
	}
}

func TestSearchEndpointWithOrFilters(
	t *testing.T,
) {

	ensureDocumentsIndex(t)

	db := createDB(t)

	cfg := createTestConfig()

	searchClient :=
		opensearch.NewClient(cfg)

	apiKeyRepository :=
		apikeys.NewRepository(db)

	apiKeyService :=
		apikeys.NewService(
			apiKeyRepository,
		)

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

	t.Cleanup(func() {

		_ = opensearch.DeleteDocument(
			context.Background(),
			searchClient,
			"search-test",
			"filter-1",
		)

		_ = opensearch.DeleteDocument(
			context.Background(),
			searchClient,
			"search-test",
			"filter-2",
		)

		_ = filterRepository.DeleteByDocumentKeys(
			context.Background(),
			[]string{
				"search-test:filter-1",
				"search-test:filter-2",
			},
		)

		_ = documentService.DeleteByNamespace(
			context.Background(),
			"search-test",
		)

		_ = apiKeyService.DeleteByNamespace(
			context.Background(),
			"search-test",
		)

		_ = db.Close()
	})

	apiKey, err :=
		apiKeyService.Create(
			t.Context(),
			"Search Test",
			"search-test",
		)

	if err != nil {
		t.Fatal(err)
	}

	err = documentService.Upsert(
		t.Context(),
		&documents.Document{
			Namespace:  "search-test",
			ExternalID: "filter-1",
			Title:      "Magento",
			Text:       "Magento ecommerce",
			Payload: []byte(`{
				"sku":"abc",
				"brand":"Dell"
			}`),
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	err = documentService.Upsert(
		t.Context(),
		&documents.Document{
			Namespace:  "search-test",
			ExternalID: "filter-2",
			Title:      "Magento",
			Text:       "Magento ecommerce",
			Payload: []byte(`{
				"sku":"xyz",
				"brand":"HP"
			}`),
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	syncDocuments(
		t,
		db,
	)

	documentHandler :=
		documents.NewHandler(
			documentService,
		)

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
			"query":"Magento",
			"filter_type":"or",
			"filters":[
				{
					"field":"sku",
					"operator":"eq",
					"value":"abc",
					"value_type":"string"
				},
				{
					"field":"brand",
					"operator":"eq",
					"value":"HP",
					"value_type":"string"
				}
			]
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

	if recorder.Code != http.StatusOK {

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

	if len(response) != 2 {

		t.Fatalf(
			"expected 2 results got %d",
			len(response),
		)
	}
}