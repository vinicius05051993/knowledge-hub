package server

import (
	"encoding/json"
	"net/http"

	"indexer/internal/auth"
)

func ProtectedHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	authContext := auth.GetAuthContext(r)

	if authContext == nil {

		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)

		return
	}

	response := map[string]any{
		"api_key_id": authContext.APIKeyID,
		"namespace": authContext.Namespace,
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		response,
	)
}