package integration

import (
	"context"
	"testing"
	"time"

	"indexer/internal/documents"
)

func TestDocumentRepository(
	t *testing.T,
) {

	db := createDB(t)

	repository :=
		documents.NewRepository(db)

	now := time.Now().UTC()

	t.Cleanup(func() {

		_ = repository.DeleteByNamespace(
			context.Background(),
			"repository-test",
		)

		_ = repository.DeleteByDocumentKeys(
			context.Background(),
			[]string{
				"repository-test:doc-1",
				"repository-test:doc-2",
			},
		)

		_ = db.Close()
	})

	document1 := &documents.Document{
		DocumentKey: "repository-test:doc-1",
		Namespace:   "repository-test",
		ExternalID:  "doc-1",
		Title:       "Document 1",
		Text:        "Document 1",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	document2 := &documents.Document{
		DocumentKey: "repository-test:doc-2",
		Namespace:   "repository-test",
		ExternalID:  "doc-2",
		Title:       "Document 2",
		Text:        "Document 2",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err := repository.Upsert(
		context.Background(),
		document1,
	)

	if err != nil {
		t.Fatal(err)
	}

	err = repository.Upsert(
		context.Background(),
		document2,
	)

	if err != nil {
		t.Fatal(err)
	}

	err = repository.MarkSyncedByDocumentKeys(
		context.Background(),
		[]string{
			"repository-test:doc-1",
			"repository-test:doc-2",
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	doc1, err :=
		repository.FindByExternalID(
			context.Background(),
			"repository-test",
			"doc-1",
		)

	if err != nil {
		t.Fatal(err)
	}

	doc2, err :=
		repository.FindByExternalID(
			context.Background(),
			"repository-test",
			"doc-2",
		)

	if err != nil {
		t.Fatal(err)
	}

	if doc1.SyncStatus !=
		documents.SyncStatusSynced {

		t.Fatalf(
			"expected sync status %d got %d",
			documents.SyncStatusSynced,
			doc1.SyncStatus,
		)
	}

	if doc2.SyncStatus !=
		documents.SyncStatusSynced {

		t.Fatalf(
			"expected sync status %d got %d",
			documents.SyncStatusSynced,
			doc2.SyncStatus,
		)
	}

	err = repository.DeleteByDocumentKeys(
		context.Background(),
		[]string{
			"repository-test:doc-1",
			"repository-test:doc-2",
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	_, err =
		repository.FindByExternalID(
			context.Background(),
			"repository-test",
			"doc-1",
		)

	if err == nil {

		t.Fatal(
			"document doc-1 should have been deleted",
		)
	}

	_, err =
		repository.FindByExternalID(
			context.Background(),
			"repository-test",
			"doc-2",
		)

	if err == nil {

		t.Fatal(
			"document doc-2 should have been deleted",
		)
	}
}