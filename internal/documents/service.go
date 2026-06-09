package documents

import (
	"context"
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
			Namespace: document.Namespace,
			ExternalID: document.ExternalID,
			Title:      document.Title,
			Text:       document.Text,
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