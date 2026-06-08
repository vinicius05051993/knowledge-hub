package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"indexer/internal/server"
)

func TestHealth(
	t *testing.T,
) {

	req := httptest.NewRequest(
		http.MethodGet,
		"/health",
		nil,
	)

	recorder := httptest.NewRecorder()

	server.HealthHandler(
		recorder,
		req,
	)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected 200 got %d",
			recorder.Code,
		)
	}
}