package middleware

import (
	"context"
	"net/http"

	"indexer/internal/apikeys"
	"indexer/internal/auth"
)

func AuthMiddleware(
	apiKeyService *apikeys.Service,
	next http.Handler,
) http.Handler {

	return http.HandlerFunc(
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {

			apiKey := r.Header.Get(
				"X-API-Key",
			)

			if apiKey == "" {
				http.Error(
					w,
					"missing api key",
					http.StatusUnauthorized,
				)

				return
			}

			key, err := apiKeyService.Validate(
				r.Context(),
				apiKey,
			)

			if err != nil {

				http.Error(
					w,
					"invalid api key",
					http.StatusUnauthorized,
				)

				return
			}

			authContext := &auth.AuthContext{
				Namespace: key.Namespace,
				APIKeyID: key.ID,
			}

			ctx := context.WithValue(
				r.Context(),
				auth.ContextKey,
				authContext,
			)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		},
	)
}