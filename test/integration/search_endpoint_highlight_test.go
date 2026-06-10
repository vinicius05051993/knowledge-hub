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

func TestSearchEndpointHighlight(
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
			"highlight-test",
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
			"Highlight Test",
		)

	if err != nil {
		t.Fatal(err)
	}

	err = documentService.Upsert(
		t.Context(),
		&documents.Document{
			Namespace:  "test",
			ExternalID: "highlight-test",
			Title:      "Magento 2.4",
			Text:       "Magento ecommerce platform",
			Payload: []byte(`{
				"sku":"highlight-test"
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
			"query":"Magento"
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

	if len(response) == 0 {

		t.Fatal(
			"expected at least one result",
		)
	}

	highlights :=
		response[0].Highlights

	if highlights == nil {

		t.Fatal(
			"expected highlights",
		)
	}

	titleHighlight :=
		highlights["title"]

	textHighlight :=
		highlights["text"]

	if titleHighlight == "" &&
		textHighlight == "" {

		t.Fatal(
			"expected title or text highlight",
		)
	}

	if !strings.Contains(
		titleHighlight+textHighlight,
		"<mark>Magento</mark>",
	) {

		t.Fatalf(
			"expected highlight containing <mark>Magento</mark>, got title=%q text=%q",
			titleHighlight,
			textHighlight,
		)
	}
}