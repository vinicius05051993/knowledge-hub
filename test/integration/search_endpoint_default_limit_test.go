package integration

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"indexer/internal/server"
)

func TestSearchDefaultLimit(
	t *testing.T,
) {

	db := createDB(t)

	defer db.Close()

	apiKey :=
		createTestAPIKey(
			t,
			db,
		)

	app :=
		server.NewApp(
			db,
			createAPIKeyService(
				t,
				db,
			),
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