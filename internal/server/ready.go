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

			json.NewEncoder(w).Encode(
				map[string]string{
					"status": "not_ready",
				},
			)

			return
		}

		json.NewEncoder(w).Encode(
			map[string]string{
				"status": "ready",
			},
		)
	}
}