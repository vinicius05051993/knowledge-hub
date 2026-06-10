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

func TestSearchEndpointPagination(
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
			"page-1",
		)

		_ = opensearch.DeleteDocument(
			context.Background(),
			searchClient,
			"test",
			"page-2",
		)

		_ = opensearch.DeleteDocument(
			context.Background(),
			searchClient,
			"test",
			"page-3",
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
			"Pagination Test",
		)

	if err != nil {
		t.Fatal(err)
	}

	documentsToCreate := []*documents.Document{
		{
			Namespace:  "test",
			ExternalID: "page-1",
			Title:      "Magento",
			Text:       "Magento ecommerce",
			Payload:    []byte(`{"sku":"pagination-test"}`),
		},
		{
			Namespace:  "test",
			ExternalID: "page-2",
			Title:      "Magento",
			Text:       "Magento ecommerce",
			Payload:    []byte(`{"sku":"pagination-test"}`),
		},
		{
			Namespace:  "test",
			ExternalID: "page-3",
			Title:      "Magento",
			Text:       "Magento ecommerce",
			Payload:    []byte(`{"sku":"pagination-test"}`),
		},
	}

	for _, document := range documentsToCreate {

		err = documentService.Upsert(
			t.Context(),
			document,
		)

		if err != nil {
			t.Fatal(err)
		}
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
			"query":"Magento",
			"offset":0,
			"limit":1
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

	var page1 []documents.SearchResponse

	err = json.Unmarshal(
		recorder.Body.Bytes(),
		&page1,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(page1) != 1 {

		t.Fatalf(
			"expected 1 result got %d",
			len(page1),
		)
	}

	req = httptest.NewRequest(
		http.MethodPost,
		"/search",
		strings.NewReader(`{
			"query":"Magento",
			"offset":1,
			"limit":1
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

	recorder =
		httptest.NewRecorder()

	router.ServeHTTP(
		recorder,
		req,
	)

	var page2 []documents.SearchResponse

	err = json.Unmarshal(
		recorder.Body.Bytes(),
		&page2,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(page2) != 1 {

		t.Fatalf(
			"expected 1 result got %d",
			len(page2),
		)
	}

	if page1[0].ExternalID ==
		page2[0].ExternalID {

		t.Fatal(
			"pagination returned same document",
		)
	}
}