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

func TestSearchEndpointPaginationWithoutQuery(
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
		for i := 1; i <= 3; i++ {
			_ = opensearch.DeleteDocument(
				context.Background(),
				searchClient,
				"test",
				string(rune('0' + i)),
			)
		}
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
			"Pagination Without Query",
		)

	if err != nil {
		t.Fatal(err)
	}

	for i := 1; i <= 3; i++ {

		err = documentService.Upsert(
			t.Context(),
			&documents.Document{
				Namespace:  "test",
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
			"filters":{
				"sku":"pagination-no-query"
			},
			"offset":0,
			"limit":1
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

	var page1 []documents.SearchResponse

	json.Unmarshal(
		rec1.Body.Bytes(),
		&page1,
	)

	req2 := httptest.NewRequest(
		http.MethodPost,
		"/search",
		strings.NewReader(`{
			"filters":{
				"sku":"pagination-no-query"
			},
			"offset":1,
			"limit":1
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

	var page2 []documents.SearchResponse

	json.Unmarshal(
		rec2.Body.Bytes(),
		&page2,
	)

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