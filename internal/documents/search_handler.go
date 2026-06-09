package documents

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) Search(
	w http.ResponseWriter,
	r *http.Request,
) {

	var request SearchRequest

	err := json.NewDecoder(
		r.Body,
	).Decode(
		&request,
	)

	if err != nil {

		http.Error(
			w,
			"invalid request",
			http.StatusBadRequest,
		)

		return
	}

	if request.Limit <= 0 {

		request.Limit = 10
	}

	results, err :=
		h.service.Search(
			r.Context(),
			request.Query,
			request.Limit,
			request.Filters,
		)

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	responses :=
		make(
			[]SearchResponse,
			0,
			len(results),
		)

	for _, document := range results {

		responses = append(
			responses,
			SearchResponse{
				Namespace: document.Namespace,

				ExternalID: document.ExternalID,

				Title: document.Title,

				Text: document.Text,

				Payload: json.RawMessage(
					document.Payload,
				),
			},
		)
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		responses,
	)
}