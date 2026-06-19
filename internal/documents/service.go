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

	err = s.filterRepository.DeleteByDocumentKeys(
		ctx,
		documentKeys,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Search(
	ctx context.Context,
	query string,
	offset int,
	limit int,
	filters map[string]string,
	filterType string,
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

	order :=
		make(
			map[string]int,
			len(results),
		)

	highlights :=
		make(
			map[string]map[string]string,
		)

	seen :=
	make(
		map[string]struct{},
	)

	for i, result := range results {

		if _, ok :=
			seen[result.DocumentKey]; ok {

			continue
		}

		seen[result.DocumentKey] =
			struct{}{}

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
			0,
			0,
			filters,
			filterType,
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

func (s *Service) prepareDocument(
	document *Document,
) error {

	if strings.TrimSpace(document.Title) == "" &&
		strings.TrimSpace(document.Text) == "" &&
		len(document.Payload) == 0 {

		return ErrEmptyDocument
	}

	document.DeletedAt = nil

	hasSearchContent :=
		strings.TrimSpace(document.Title) != "" ||
			strings.TrimSpace(document.Text) != ""

	if hasSearchContent {
		document.SyncStatus =
			SyncStatusPendingUpsert
	} else {
		document.SyncStatus =
			SyncStatusSynced
	}

	now := time.Now().UTC()

	if document.CreatedAt.IsZero() {
		document.CreatedAt = now
	}

	document.UpdatedAt = now

	document.DocumentKey =
		document.Namespace +
			":" +
			document.ExternalID

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
			map[string]map[string]string,
		)

	for _, document := range documents {

		filters :=
			make(
				map[string]string,
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

				str, ok := value.(string)

				if !ok {
					continue
				}

				filters[key] = str
			}
		}

		filtersByDocument[document.DocumentKey] = filters
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