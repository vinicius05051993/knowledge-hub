package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"indexer/internal/apikeys"
	"indexer/internal/server"
	"indexer/test/api/mocks"
)

func TestProtectedInvalidAPIKey(
	t *testing.T,
) {

	repository :=
		&mocks.APIKeyRepository{
			FindByHashFunc: func(
				ctx context.Context,
				hash string,
			) (*apikeys.APIKey, error) {

				return nil,
					errors.New("not found")
			},
		}

	service :=
		apikeys.NewService(
			repository,
		)

	app :=
		server.NewApp(
			nil,
			service,
		)

	router := app.Router()

	req := httptest.NewRequest(
		http.MethodGet,
		"/protected",
		nil,
	)

	req.Header.Set(
		"X-API-Key",
		"invalid-key",
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(
		recorder,
		req,
	)

	if recorder.Code != http.StatusUnauthorized {

		t.Fatalf(
			"expected 401 got %d",
			recorder.Code,
		)
	}
}

func TestProtectedValidAPIKey(
	t *testing.T,
) {

	repository :=
		&mocks.APIKeyRepository{
			FindByHashFunc: func(
				ctx context.Context,
				hash string,
			) (*apikeys.APIKey, error) {

				return &apikeys.APIKey{
					ID:        1,
					Namespace: "magento",
				}, nil
			},
		}

	service :=
		apikeys.NewService(
			repository,
		)

	app :=
		server.NewApp(
			nil,
			service,
		)

	router := app.Router()

	req := httptest.NewRequest(
		http.MethodGet,
		"/protected",
		nil,
	)

	req.Header.Set(
		"X-API-Key",
		"sk_live_test",
	)

	recorder := httptest.NewRecorder()

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
}