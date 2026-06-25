package integration

import (
	"context"
	"testing"
	"time"

	"indexer/internal/documents"
	"indexer/internal/opensearch"
	"indexer/internal/documentfilters"
)

func TestSyncServiceProcessPendingDeletes(
	t *testing.T,
) {

	ensureDocumentsIndex(t)

	db := createDB(t)

	cfg := createTestConfig()

	searchClient :=
		opensearch.NewClient(cfg)

	searchService :=
		opensearch.NewService(
			searchClient,
		)

	documentRepository :=
		documents.NewRepository(db)

	t.Cleanup(func() {

		_ = opensearch.DeleteDocument(
			context.Background(),
			searchClient,
			"sync-test",
			"test-delete",
		)

		_ = documentRepository.DeleteByNamespace(
			context.Background(),
			"sync-test",
		)

		syncDocuments(
			t,
			db,
		)

		_ = db.Close()
	})

	filterRepository :=
		documentfilters.NewRepository(
			db,
		)

	syncService :=
		documents.NewSyncService(
			documentRepository,
			filterRepository,
			searchService,
		)

	now := time.Now().UTC()

	document := &documents.Document{
		DocumentKey: "sync-test:test-delete",

		Namespace:  "sync-test",
		ExternalID: "test-delete",

		Title: "Delete Test",
		Text:  "Document To Delete",

		SyncStatus:
			documents.SyncStatusPendingUpsert,

		CreatedAt: now,
		UpdatedAt: now,
	}

	err := documentRepository.Upsert(
		context.Background(),
		document,
	)

	if err != nil {
		t.Fatal(err)
	}

	err = syncService.ProcessPendingUpserts(
		context.Background(),
		100,
	)

	if err != nil {
		t.Fatal(err)
	}

	err = documentRepository.DeleteByExternalIDs(
		context.Background(),
		"sync-test",
		[]string{
			"test-delete",
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	err = syncService.ProcessPendingDeletes(
		context.Background(),
		100,
	)

	if err != nil {
		t.Fatal(err)
	}

	_, err =
		documentRepository.FindByExternalID(
			context.Background(),
			"sync-test",
			"test-delete",
		)

	if err == nil {

		t.Fatal(
			"document should have been deleted",
		)
	}

	results, err :=
		searchService.Search(
			context.Background(),
			"Delete",
			0,
			10,
		)

	if err != nil {
		t.Fatal(err)
	}

	for _, result := range results {

		if result.DocumentKey ==
			"sync-test:test-delete" {

			t.Fatal(
				"document still exists in opensearch",
			)
		}
	}
}