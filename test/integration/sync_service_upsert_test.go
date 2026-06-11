package integration

import (
	"context"
	"testing"
	"time"

	"indexer/internal/documents"
	"indexer/internal/opensearch"
)

func TestSyncServiceProcessPendingUpserts(
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
			"test-1",
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

	syncService :=
		documents.NewSyncService(
			documentRepository,
			searchService,
		)

	now := time.Now().UTC()

	document := &documents.Document{
		DocumentKey: "sync-test:test-1",

		Namespace:  "sync-test",
		ExternalID: "test-1",

		Title: "Adobe Experience Manager",
		Text:  "Adobe AEM Search Test",

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

	documentDB, err :=
		documentRepository.FindByExternalID(
			context.Background(),
			"sync-test",
			"test-1",
		)

	if err != nil {
		t.Fatal(err)
	}

	if documentDB.SyncStatus !=
		documents.SyncStatusSynced {

		t.Fatalf(
			"expected sync_status=%d got=%d",
			documents.SyncStatusSynced,
			documentDB.SyncStatus,
		)
	}

	results, err :=
		searchService.Search(
			context.Background(),
			"AEM",
			0,
			10,
		)

	if err != nil {
		t.Fatal(err)
	}

	found := false

	for _, result := range results {

		if result.DocumentKey ==
			"sync-test:test-1" {

			found = true
			break
		}
	}

	if !found {

		t.Fatal(
			"document not indexed",
		)
	}
}