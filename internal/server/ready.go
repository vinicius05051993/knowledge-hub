package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

func ReadyHandler(
	db *sqlx.DB,
) http.HandlerFunc {

	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		w.Header().Set(
			"Content-Type",
			"application/json",
		)

		ctx, cancel := context.WithTimeout(
			r.Context(),
			2*time.Second,
		)

		defer cancel()

		err := db.PingContext(ctx)

		if err != nil {

			w.WriteHeader(
				http.StatusServiceUnavailable,
			)

			_ = json.NewEncoder(w).Encode(
				map[string]string{
					"status": "not_ready",
					"mysql":  "error",
				},
			)

			return
		}

		w.WriteHeader(
			http.StatusOK,
		)

		_ = json.NewEncoder(w).Encode(
			map[string]string{
				"status": "ready",
				"mysql":  "ok",
			},
		)
	}
}