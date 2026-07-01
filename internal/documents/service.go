package documents

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"time"

	"indexer/internal/documentfilters"
	"indexer/internal/opensearch"
)

const (
	FilterTypeAnd = "and"
	FilterTypeOr  = "or"
)

var ErrEmptyDocument = errors.New(
	"document must contain title, text or payload",
)

type SearchIndexer interface {
	Search(
		ctx context.Context,
		query string,
		offset int,
		limit int,
	) ([]opensearch.SearchResult, error)
}

type Service struct {
	repository       *Repository
	filterRepository *documentfilters.Repository
	indexer          SearchIndexer
}

func NewService(
	repository *Repository,
	filterRepository *documentfilters.Repository,
	indexer SearchIndexer,
) *Service {

	return &Service{
		repository:       repository,
		filterRepository: filterRepository,
		indexer:          indexer,
	}
}

func (s *Service) DeleteByNamespace(
	ctx context.Context,
	namespace string,
) error {

	return s.repository.DeleteByNamespace(
		ctx,
		namespace,
	)
}

func (s *Service) Upsert(
	ctx context.Context,
	document *Document,
) error {

	return s.UpsertBatch(
		ctx,
		[]*Document{
			document,
		},
	)
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

	documentKeys :=
		make(
			[]string,
			0,
			len(externalIDs),
		)

	for _, externalID := range externalIDs {

		documentKeys = append(
			documentKeys,
			namespace+":"+externalID,
		)
	}
	
	return nil
}

func (s *Service) Search(
	ctx context.Context,
	query string,
	offset int,
	limit int,
	filters []SearchFilter,
	filterType string,
	order *SearchOrder,
) ([]SearchDocument, error) {

	if filterType != FilterTypeOr {
		filterType = FilterTypeAnd
	}

	if query == "" {

		documents, err := s.repository.Search(
			ctx,
			nil,
			offset,
			limit,
			filters,
			filterType,
			order,
		)

		if err != nil {
			return nil, err
		}

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

	documentOrder :=
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

		documentOrder[
			result.DocumentKey,
		] = i
	}

	documents, err :=
		s.repository.Search(
			ctx,
			documentKeys,
			0,
			0,
			filters,
			filterType,
			order,
		)

	if err != nil {
		return nil, err
	}

	sort.Slice(
		documents,
		func(i, j int) bool {

			return documentOrder[
				documents[i].DocumentKey,
			] <
				documentOrder[
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

func (s *Service) prepareDocument(
	ctx context.Context,
	document *Document,
) error {

	if strings.TrimSpace(document.Title) == "" &&
		strings.TrimSpace(document.Text) == "" &&
		len(document.Payload) == 0 {

		return ErrEmptyDocument
	}

	document.DeletedAt = nil

	document.DocumentKey =
		document.Namespace +
			":" +
			document.ExternalID

	hasSearchContent :=
		strings.TrimSpace(document.Title) != "" ||
			strings.TrimSpace(document.Text) != ""

	existing, err :=
		s.repository.FindByExternalID(
			ctx,
			document.Namespace,
			document.ExternalID,
		)

	if err == nil {

		existingHasSearchContent :=
			strings.TrimSpace(existing.Title) != "" ||
				strings.TrimSpace(existing.Text) != ""

		switch {

		case hasSearchContent:

			document.SyncStatus =
				SyncStatusPendingUpsert

		case existingHasSearchContent:

			document.SyncStatus =
				SyncStatusPendingDeindex

		default:

			document.SyncStatus =
				SyncStatusSynced
		}

		document.CreatedAt =
			existing.CreatedAt

	} else {

		if hasSearchContent {

			document.SyncStatus =
				SyncStatusPendingUpsert

		} else {

			document.SyncStatus =
				SyncStatusSynced
		}
	}

	now := time.Now().UTC()

	if document.CreatedAt.IsZero() {
		document.CreatedAt = now
	}

	document.UpdatedAt = now

	return nil
}

func (s *Service) UpsertBatch(
	ctx context.Context,
	documents []*Document,
) error {

	if len(documents) == 0 {
		return nil
	}

	for _, document := range documents {

		err := s.prepareDocument(
			ctx,
			document,
		)

		if err != nil {
			return err
		}
	}

	err := s.repository.UpsertBatch(
		ctx,
		documents,
	)

	if err != nil {
		return err
	}

	filtersByDocument :=
		make(
			map[string]map[string][]string,
		)

	for _, document := range documents {

		filters := make(
			map[string][]string,
		)

		if len(document.Payload) > 0 {

			var payload map[string]any

			err = json.Unmarshal(
				document.Payload,
				&payload,
			)

			if err != nil {
				return err
			}

			for key, value := range payload {

				switch v := value.(type) {

				case string:

					filters[key] = append(
						filters[key],
						v,
					)

				case []any:

					for _, item := range v {

						str, ok := item.(string)

						if !ok {
							continue
						}

						filters[key] = append(
							filters[key],
							str,
						)
					}
				}
			}
		}

		filtersByDocument[
			document.DocumentKey,
		] = filters
	}

	err = s.filterRepository.ReplaceBatch(
		ctx,
		filtersByDocument,
	)

	if err != nil {
		return err
	}

	return nil
}