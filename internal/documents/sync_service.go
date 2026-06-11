package documents

import (
	"context"

	"indexer/internal/opensearch"
)

type SyncIndexer interface {
	IndexDocument(
		ctx context.Context,
		document *opensearch.Document,
	) error

	DeleteDocuments(
		ctx context.Context,
		documentKeys []string,
	) error
}

type SyncService struct {
	repository *Repository
	indexer    SyncIndexer
}

func NewSyncService(
	repository *Repository,
	indexer SyncIndexer,
) *SyncService {

	return &SyncService{
		repository: repository,
		indexer:    indexer,
	}
}

func (s *SyncService) ProcessPendingUpserts(
	ctx context.Context,
	limit int,
) error {

	documents, err :=
		s.repository.FindPendingUpserts(
			ctx,
			limit,
		)

	if err != nil {
		return err
	}

	for _, document := range documents {

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
			continue
		}

		err = s.repository.MarkSynced(
			ctx,
			document.DocumentKey,
		)

		if err != nil {
			continue
		}
	}

	return nil
}

func (s *SyncService) ProcessPendingDeletes(
	ctx context.Context,
	limit int,
) error {

	documents, err :=
		s.repository.FindPendingDeletes(
			ctx,
			limit,
		)

	if err != nil {
		return err
	}

	for _, document := range documents {

		err = s.indexer.DeleteDocuments(
			ctx,
			[]string{
				document.DocumentKey,
			},
		)

		if err != nil {
			continue
		}

		err = s.repository.DeleteByDocumentKey(
			ctx,
			document.DocumentKey,
		)

		if err != nil {
			continue
		}
	}

	return nil
}