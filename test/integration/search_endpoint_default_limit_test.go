package integration

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"context"

	"indexer/internal/apikeys"
	"indexer/internal/server"
)

func TestSearchDefaultLimit(
	t *testing.T,
) {

	db := createDB(t)

	apiKeyRepository :=
		apikeys.NewRepository(db)

	apiKeyService :=
		apikeys.NewService(
			apiKeyRepository,
		)

	t.Cleanup(func() {

		_ = apiKeyService.DeleteByNamespace(
			context.Background(),
			"search-default-limit-test",
		)

		_ = db.Close()
	})

	apiKey, err :=
		apiKeyService.Create(
			t.Context(),
			"Search Default Limit Test",
			"search-default-limit-test",
		)

	if err != nil {
		t.Fatal(err)
	}

	app :=
		server.NewApp(
			db,
			apiKeyService,
			createDocumentHandler(
				t,
				db,
			),
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
		"X-API-Key",
		apiKey,
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	rec :=
		httptest.NewRecorder()

	router.ServeHTTP(
		rec,
		req,
	)

	if rec.Code !=
		http.StatusOK {

		t.Fatalf(
			"expected 200 got %d",
			rec.Code,
		)
	}
}