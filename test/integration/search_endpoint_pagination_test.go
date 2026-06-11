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

func TestSearchEndpointPagination(
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
			"pagination-test",
			"page-1",
		)

		_ = opensearch.DeleteDocument(
			context.Background(),
			searchClient,
			"pagination-test",
			"page-2",
		)

		_ = opensearch.DeleteDocument(
			context.Background(),
			searchClient,
			"pagination-test",
			"page-3",
		)

		_ = filterRepository.DeleteByDocumentKeys(
			context.Background(),
			[]string{
				"pagination-test:page-1",
				"pagination-test:page-2",
				"pagination-test:page-3",
			},
		)

		_ = documentService.DeleteByNamespace(
			context.Background(),
			"pagination-test",
		)

		_ = apiKeyService.DeleteByNamespace(
			context.Background(),
			"pagination-test",
		)

		_ = db.Close()
	})

	documentHandler :=
		documents.NewHandler(
			documentService,
		)

	apiKey, err :=
		apiKeyService.Create(
			t.Context(),
			"Pagination Test",
			"pagination-test",
		)

	if err != nil {
		t.Fatal(err)
	}

	documentsToCreate := []*documents.Document{
		{
			Namespace:  "pagination-test",
			ExternalID: "page-1",
			Title:      "Magento",
			Text:       "Magento ecommerce",
			Payload:    []byte(`{"sku":"pagination-test"}`),
		},
		{
			Namespace:  "pagination-test",
			ExternalID: "page-2",
			Title:      "Magento",
			Text:       "Magento ecommerce",
			Payload:    []byte(`{"sku":"pagination-test"}`),
		},
		{
			Namespace:  "pagination-test",
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

	syncDocuments(
		t,
		db,
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

	if recorder.Code !=
		http.StatusOK {

		t.Fatalf(
			"expected 200 got %d",
			recorder.Code,
		)
	}

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