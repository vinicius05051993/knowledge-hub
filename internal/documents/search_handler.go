package documents

import (
	"encoding/json"
	"net/http"
	"time"

	"indexer/internal/metrics"
)

func (h *Handler) Search(
	w http.ResponseWriter,
	r *http.Request,
) {

	metrics.SearchRequestsTotal.Inc()

	start := time.Now()

	defer func() {
		metrics.SearchDurationSeconds.Observe(
			time.Since(start).Seconds(),
		)
	}()

	var request SearchRequest

	err := json.NewDecoder(
		r.Body,
	).Decode(
		&request,
	)

	if err != nil {

		metrics.SearchErrorsTotal.Inc()

		http.Error(
			w,
			"invalid request",
			http.StatusBadRequest,
		)

		return
	}

	if request.Offset < 0 {

		request.Offset = 0
	}

	if request.Limit <= 0 {

		request.Limit = 100
	}

	if request.Limit > 100 {

		request.Limit = 100
	}

	err = normalizeOrder(request.Order)

	if err != nil {

		metrics.SearchErrorsTotal.Inc()

		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)

		return
	}

	results, err :=
		h.service.Search(
			r.Context(),
			request.Query,
			request.Offset,
			request.Limit,
			request.Filters,
			request.FilterType,
			request.Order,
		)

	if err != nil {

		metrics.SearchErrorsTotal.Inc()

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

				Highlights: document.Highlights,
			},
		)
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	_ = json.NewEncoder(w).Encode(
		responses,
	)
}