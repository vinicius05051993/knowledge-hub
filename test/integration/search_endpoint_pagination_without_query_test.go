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

func TestSearchEndpointPaginationWithoutQuery(
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

		for i := 1; i <= 3; i++ {

			_ = opensearch.DeleteDocument(
				context.Background(),
				searchClient,
				"pwq-test",
				string(rune('0'+i)),
			)
		}

		_ = filterRepository.DeleteByDocumentKeys(
			context.Background(),
			[]string{
				"pwq-test:1",
				"pwq-test:2",
				"pwq-test:3",
			},
		)

		_ = documentService.DeleteByNamespace(
			context.Background(),
			"pwq-test",
		)

		_ = apiKeyService.DeleteByNamespace(
			context.Background(),
			"pwq-test",
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
			"Pagination Without Query",
			"pwq-test",
		)

	if err != nil {
		t.Fatal(err)
	}

	for i := 1; i <= 3; i++ {

		err = documentService.Upsert(
			t.Context(),
			&documents.Document{
				Namespace:  "pwq-test",
				ExternalID: string(rune('0' + i)),
				Title:      "Documento",
				Text:       "Texto",
				Payload:    []byte(`{"sku":"pagination-no-query"}`),
			},
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

	req1 := httptest.NewRequest(
		http.MethodPost,
		"/search",
		strings.NewReader(`{
			"filters": [
				{
					"field": "sku",
					"operator": "eq",
					"value": "pagination-no-query",
					"value_type": "string"
				}
			],
			"offset": 0,
			"limit": 1
		}`),
	)

	req1.Header.Set(
		"Content-Type",
		"application/json",
	)

	req1.Header.Set(
		"X-API-Key",
		apiKey,
	)

	rec1 :=
		httptest.NewRecorder()

	router.ServeHTTP(
		rec1,
		req1,
	)

	if rec1.Code !=
		http.StatusOK {

		t.Fatalf(
			"expected 200 got %d",
			rec1.Code,
		)
	}

	var page1 []documents.SearchResponse

	err = json.Unmarshal(
		rec1.Body.Bytes(),
		&page1,
	)

	if err != nil {
		t.Fatal(err)
	}

	req2 := httptest.NewRequest(
		http.MethodPost,
		"/search",
		strings.NewReader(`{
			"filters": [
				{
					"field": "sku",
					"operator": "eq",
					"value": "pagination-no-query",
					"value_type": "string"
				}
			],
			"offset": 1,
			"limit": 1
		}`),
	)

	req2.Header.Set(
		"Content-Type",
		"application/json",
	)

	req2.Header.Set(
		"X-API-Key",
		apiKey,
	)

	rec2 :=
		httptest.NewRecorder()

	router.ServeHTTP(
		rec2,
		req2,
	)

	if rec2.Code !=
		http.StatusOK {

		t.Fatalf(
			"expected 200 got %d",
			rec2.Code,
		)
	}

	var page2 []documents.SearchResponse

	err = json.Unmarshal(
		rec2.Body.Bytes(),
		&page2,
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