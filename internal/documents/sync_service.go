package documents

import (
	"context"
	"strings"

	"indexer/internal/documentfilters"
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
	filterRepository *documentfilters.Repository
	indexer    SyncIndexer
}

func NewSyncService(
	repository *Repository,
	filterRepository *documentfilters.Repository,
	indexer SyncIndexer,
) *SyncService {

	return &SyncService{
		repository: repository,
		filterRepository: filterRepository,
		indexer:    indexer,
	}
}

func (s *SyncService) ProcessPendingDeindexes(
	ctx context.Context,
	limit int,
) error {

	pending, err :=
		s.repository.CountPendingDeindexes(
			ctx,
		)

	if err == nil {

		metrics.PendingDeletes.Set(
			float64(pending),
		)
	}

	documents, err :=
		s.repository.FindPendingDeindexes(
			ctx,
			limit,
		)

	if err != nil {

		metrics.SyncDeleteErrorsTotal.Inc()

		return err
	}

	if len(documents) == 0 {

		metrics.PendingDeletes.Set(0)

		return nil
	}

	documentKeys :=
		make(
			[]string,
			0,
			len(documents),
		)

	for _, document := range documents {

		documentKeys = append(
			documentKeys,
			document.DocumentKey,
		)
	}

	err = s.indexer.DeleteDocuments(
		ctx,
		documentKeys,
	)

	if err != nil {

		metrics.SyncDeleteErrorsTotal.Inc()

		return err
	}

	metrics.SyncDeletesTotal.Add(
		float64(
			len(documentKeys),
		),
	)

	err = s.repository.MarkSyncedByDocumentKeys(
		ctx,
		documentKeys,
	)

	if err != nil {

		metrics.SyncDeleteErrorsTotal.Inc()

		return err
	}

	pending, err =
		s.repository.CountPendingDeindexes(
			ctx,
		)

	if err == nil {

		metrics.PendingDeletes.Set(
			float64(pending),
		)
	}

	return nil
}

func (s *SyncService) ProcessPendingUpserts(
	ctx context.Context,
	limit int,
) error {

	pending, err :=
		s.repository.CountPendingUpserts(
			ctx,
		)

	if err == nil {

		metrics.PendingUpserts.Set(
			float64(pending),
		)
	}

	documents, err :=
		s.repository.FindPendingUpserts(
			ctx,
			limit,
		)

	if err != nil {

		metrics.SyncUpsertErrorsTotal.Inc()

		return err
	}

	if len(documents) == 0 {

		metrics.PendingUpserts.Set(0)

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

			metrics.SyncUpsertErrorsTotal.Inc()

			return err
		}

		metrics.SyncUpsertsTotal.Add(
			float64(
				len(bulkDocuments),
			),
		)
	}

	err = s.repository.MarkSyncedByDocumentKeys(
		ctx,
		documentKeys,
	)

	if err != nil {

		metrics.SyncUpsertErrorsTotal.Inc()

		return err
	}

	pending, err =
		s.repository.CountPendingUpserts(
			ctx,
		)

	if err == nil {

		metrics.PendingUpserts.Set(
			float64(pending),
		)
	}

	return nil
}

func (s *SyncService) ProcessPendingDeletes(
	ctx context.Context,
	limit int,
) error {

	pending, err :=
		s.repository.CountPendingDeletes(
			ctx,
		)

	if err == nil {

		metrics.PendingDeletes.Set(
			float64(pending),
		)
	}

	documents, err :=
		s.repository.FindPendingDeletes(
			ctx,
			limit,
		)

	if err != nil {

		metrics.SyncDeleteErrorsTotal.Inc()

		return err
	}

	if len(documents) == 0 {

		metrics.PendingDeletes.Set(0)

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

		metrics.SyncDeleteErrorsTotal.Inc()

		return err
	}

	metrics.SyncDeletesTotal.Add(
		float64(len(documentKeys)),
	)

	err = s.filterRepository.DeleteByDocumentKeys(
		ctx,
		documentKeys,
	)

	if err != nil {

		metrics.SyncDeleteErrorsTotal.Inc()

		return err
	}

	err = s.repository.DeleteByDocumentKeys(
		ctx,
		documentKeys,
	)

	if err != nil {

		metrics.SyncDeleteErrorsTotal.Inc()

		return err
	}

	pending, err =
		s.repository.CountPendingDeletes(
			ctx,
		)

	if err == nil {

		metrics.PendingDeletes.Set(
			float64(pending),
		)
	}

	return nil
}