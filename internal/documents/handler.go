package documents

import (
	"context"
	"encoding/json"
	"net/http"

	"indexer/internal/auth"
)

type SearchService interface {
	Search(
		ctx context.Context,
		query string,
		offset int,
		limit int,
		filters map[string]string,
	) ([]SearchDocument, error)

	Upsert(
		ctx context.Context,
		document *Document,
	) error

	Delete(
		ctx context.Context,
		namespace string,
		externalIDs []string,
	) error
}

type Handler struct {
	service SearchService
}

func NewHandler(
	service *Service,
) *Handler {

	return &Handler{
		service: service,
	}
}

func (h *Handler) UpsertDocuments(
	w http.ResponseWriter,
	r *http.Request,
) {

	authContext :=
		auth.GetAuthContext(r)

	if authContext == nil {

		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)

		return
	}

	var request UpsertDocumentsRequest

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

	for _, item := range request.Documents {

		var payload []byte

		if item.Payload != nil {

			payload, err = json.Marshal(
				item.Payload,
			)

			if err != nil {

				http.Error(
					w,
					"invalid payload",
					http.StatusBadRequest,
				)

				return
			}
		}

		document := &Document{
			Namespace: authContext.Namespace,
			ExternalID: item.ExternalID,
			Title: item.Title,
			Text: item.Text,
			Payload: payload,
		}

		err = h.service.Upsert(
			context.Background(),
			document,
		)

		if err != nil {

			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)

			return
		}
	}

	response := UpsertDocumentsResponse{
		Success: true,
		Processed: len(
			request.Documents,
		),
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	_ = json.NewEncoder(w).Encode(
		response,
	)
}

func (h *Handler) DeleteDocuments(
	w http.ResponseWriter,
	r *http.Request,
) {

	authContext :=
		auth.GetAuthContext(r)

	if authContext == nil {

		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)

		return
	}

	var request DeleteDocumentsRequest

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

	err = h.service.Delete(
		r.Context(),
		authContext.Namespace,
		request.ExternalIDs,
	)

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	w.WriteHeader(
		http.StatusNoContent,
	)
}