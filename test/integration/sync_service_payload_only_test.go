package integration

import (
	"context"
	"testing"
	"time"

	"indexer/internal/documents"
	"indexer/internal/opensearch"
	"indexer/internal/documentfilters"
)

func TestSyncServiceProcessPendingUpsertsPayloadOnly(
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

		_ = documentRepository.DeleteByNamespace(
			context.Background(),
			"payload-sync",
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
		DocumentKey:
			"payload-sync:payload-only",

		Namespace: "payload-sync",

		ExternalID: "payload-only",

		SyncStatus:
			documents.SyncStatusSynced,

		Payload: []byte(`{
			"sku":"ABC123"
		}`),

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
			"payload-sync",
			"payload-only",
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
			"ABC123",
			0,
			10,
		)

	if err != nil {
		t.Fatal(err)
	}

	for _, result := range results {

		if result.DocumentKey ==
			"payload-sync:payload-only" {

			t.Fatal(
				"payload only document should not be indexed",
			)
		}
	}
}