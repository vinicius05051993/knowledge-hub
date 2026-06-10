package documents

import (
	"context"
	"sort"
	"time"

	"indexer/internal/opensearch"
)

type SearchIndexer interface {

	IndexDocument(
		ctx context.Context,
		document *opensearch.Document,
	) error

	DeleteDocument(
		ctx context.Context,
		namespace string,
		externalID string,
	) error

	Search(
		ctx context.Context,
		query string,
		offset int,
		limit int,
	) ([]opensearch.SearchResult, error)
}

type Service struct {
	repository *Repository
	indexer    SearchIndexer
}

func NewService(
	repository *Repository,
	indexer SearchIndexer,
) *Service {

	return &Service{
		repository: repository,
		indexer:    indexer,
	}
}

func (s *Service) Upsert(
	ctx context.Context,
	document *Document,
) error {

	now := time.Now().UTC()

	if document.CreatedAt.IsZero() {
		document.CreatedAt = now
	}

	document.UpdatedAt = now

	document.DocumentKey =
		document.Namespace +
			":" +
			document.ExternalID

	err := s.repository.Upsert(
		ctx,
		document,
	)

	if err != nil {
		return err
	}

	err = s.indexer.IndexDocument(
		ctx,
		&opensearch.Document{
			DocumentKey: document.DocumentKey,
			Namespace:   document.Namespace,
			ExternalID:  document.ExternalID,
			Title:       document.Title,
			Text:        document.Text,
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Delete(
	ctx context.Context,
	namespace string,
	externalIDs []string,
) error {

	err := s.repository.DeleteByExternalIDs(
		ctx,
		namespace,
		externalIDs,
	)

	if err != nil {
		return err
	}

	for _, externalID := range externalIDs {

		err = s.indexer.DeleteDocument(
			ctx,
			namespace,
			externalID,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) Search(
	ctx context.Context,
	query string,
	offset int,
	limit int,
	filters map[string]string,
) ([]SearchDocument, error) {

	if query == "" {

		documents, err := s.repository.Search(
			ctx,
			nil,
			filters,
			offset,
			limit,
		)

		if err != nil {
			return nil, err
		}

		response := make(
			[]SearchDocument,
			0,
			len(documents),
		)

		for _, document := range documents {

			response = append(
				response,
				SearchDocument{
					Document: document,
				},
			)
		}

		return response, nil
	}

	results, err :=
		s.indexer.Search(
			ctx,
			query,
			offset,
			limit,
		)

	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return []SearchDocument{}, nil
	}

	documentKeys :=
		make(
			[]string,
			0,
			len(results),
		)

	order :=
		make(
			map[string]int,
			len(results),
		)

	highlights :=
		make(
			map[string]map[string]string,
		)

	for i, result := range results {

		normalizedHighlights :=
			make(
				map[string]string,
			)

		for field, values :=
			range result.Highlights {

			if len(values) == 0 {
				continue
			}

			normalizedHighlights[field] =
				values[0]
		}

		highlights[
			result.DocumentKey,
		] = normalizedHighlights

		documentKeys = append(
			documentKeys,
			result.DocumentKey,
		)

		order[
			result.DocumentKey,
		] = i
	}

	documents, err :=
		s.repository.Search(
			ctx,
			documentKeys,
			filters,
			0,
			0,
		)

	if err != nil {
		return nil, err
	}

	sort.Slice(
		documents,
		func(i, j int) bool {

			return order[
				documents[i].DocumentKey,
			] <
				order[
					documents[j].DocumentKey,
				]
		},
	)

	response :=
		make(
			[]SearchDocument,
			0,
			len(documents),
		)

	for _, document := range documents {

		response = append(
			response,
			SearchDocument{
				Document: document,

				Highlights: highlights[
					document.DocumentKey,
				],
			},
		)
	}

	return response, nil
}