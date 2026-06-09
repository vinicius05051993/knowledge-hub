package integration

import (
	"context"
	"testing"
	"time"

	"indexer/internal/documents"
)

func TestFindByDocumentKeys(
	t *testing.T,
) {

	db := createDB(t)

	defer db.Close()

	cleanupTestData(
		t,
		db,
	)

	repository :=
		documents.NewRepository(db)

	now := time.Now().UTC()

	document1 := &documents.Document{
		DocumentKey: "test:1",

		Namespace: "test",

		ExternalID: "1",

		Title: "Document 1",

		Text: "Text 1",

		Payload: []byte(`{
			"id":1
		}`),

		CreatedAt: now,

		UpdatedAt: now,
	}

	document2 := &documents.Document{
		DocumentKey: "test:2",

		Namespace: "test",

		ExternalID: "2",

		Title: "Document 2",

		Text: "Text 2",

		Payload: []byte(`{
			"id":2
		}`),

		CreatedAt: now,

		UpdatedAt: now,
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

	foundDocuments, err :=
		repository.FindByDocumentKeys(
			context.Background(),
			[]string{
				"test:1",
				"test:2",
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	if len(foundDocuments) != 2 {

		t.Fatalf(
			"expected 2 documents, got %d",
			len(foundDocuments),
		)
	}
}