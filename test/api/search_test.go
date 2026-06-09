package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"indexer/internal/server"
)

func TestSearchWithoutAPIKey(
	t *testing.T,
) {

	app :=
		server.NewApp(
			nil,
			nil,
			nil,
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

	recorder :=
		httptest.NewRecorder()

	router.ServeHTTP(
		recorder,
		req,
	)

	if recorder.Code !=
		http.StatusUnauthorized {

		t.Fatalf(
			"expected 401 got %d",
			recorder.Code,
		)
	}
}