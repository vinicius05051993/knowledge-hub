package documents

import (
	"context"
	"strings"

	"indexer/internal/metrics"
	"indexer/internal/opensearch"
)

type SyncIndexer interface {
	IndexDocument(
		ctx context.Context,
		document *opensearch.Document,
	) error

	BulkIndexDocuments(
		ctx context.Context,
		documents []*opensearch.Document,
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

	if len(documents) == 0 {
		return nil
	}

	bulkDocuments :=
		make(
			[]*opensearch.Document,
			0,
			len(documents),
		)

	documentKeys :=
		make(
			[]string,
			0,
			len(documents),
		)

	for _, document := range documents {

		documentKeys =
			append(
				documentKeys,
				document.DocumentKey,
			)

		// Documentos somente com payload
		// ficam apenas no MySQL
		if strings.TrimSpace(
			document.Title,
		) == "" &&
			strings.TrimSpace(
				document.Text,
			) == "" {

			continue
		}

		bulkDocuments =
			append(
				bulkDocuments,
				&opensearch.Document{
					DocumentKey: document.DocumentKey,
					Namespace:   document.Namespace,
					ExternalID:  document.ExternalID,
					Title:       document.Title,
					Text:        document.Text,
				},
			)
	}

	if len(bulkDocuments) > 0 {

		err = s.indexer.BulkIndexDocuments(
			ctx,
			bulkDocuments,
		)

		if err != nil {
			return err
		}

		metrics.SyncUpsertsTotal.Add(
			float64(
				len(bulkDocuments),
			),
		)
	}

	return s.repository.MarkSyncedByDocumentKeys(
		ctx,
		documentKeys,
	)
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

	if len(documents) == 0 {
		return nil
	}

	documentKeys :=
		make(
			[]string,
			0,
			len(documents),
		)

	for _, document := range documents {

		documentKeys =
			append(
				documentKeys,
				document.DocumentKey,
			)
	}

	err = s.indexer.DeleteDocuments(
		ctx,
		documentKeys,
	)

	if err != nil {
		return err
	}

	metrics.SyncDeletesTotal.Add(
		float64(
			len(documentKeys),
		),
	)

	return s.repository.DeleteByDocumentKeys(
		ctx,
		documentKeys,
	)
}