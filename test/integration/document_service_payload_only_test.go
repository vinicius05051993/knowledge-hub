package integration

import (
	"context"
	"testing"

	"indexer/internal/documentfilters"
	"indexer/internal/documents"
	"indexer/internal/opensearch"
)

func TestDocumentServiceUpsertPayloadOnly(
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

	filterRepository :=
		documentfilters.NewRepository(
			db,
		)

	documentService :=
		documents.NewService(
			documentRepository,
			filterRepository,
			searchService,
		)

	t.Cleanup(func() {

		_ = documentService.DeleteByNamespace(
			context.Background(),
			"payload-test",
		)

		_ = filterRepository.DeleteByDocumentKeys(
			context.Background(),
			[]string{
				"payload-test:payload-only",
			},
		)

		syncDocuments(
			t,
			db,
		)

		_ = db.Close()
	})

	document := &documents.Document{
		Namespace:  "payload-test",
		ExternalID: "payload-only",
		Payload: []byte(`{
			"sku":"ABC123"
		}`),
	}

	err := documentService.Upsert(
		context.Background(),
		document,
	)

	if err != nil {
		t.Fatal(err)
	}

	foundDocument, err :=
		documentRepository.FindByExternalID(
			context.Background(),
			"payload-test",
			"payload-only",
		)

	if err != nil {
		t.Fatal(err)
	}

	if foundDocument == nil {

		t.Fatal(
			"document not found",
		)
	}
}

func TestDocumentServiceUpsertEmptyDocument(
	t *testing.T,
) {

	db := createDB(t)

	documentRepository :=
		documents.NewRepository(db)

	filterRepository :=
		documentfilters.NewRepository(
			db,
		)

	documentService :=
		documents.NewService(
			documentRepository,
			filterRepository,
			nil,
		)

	t.Cleanup(func() {
		_ = db.Close()
	})

	document := &documents.Document{
		Namespace:  "empty-test",
		ExternalID: "empty",
	}

	err := documentService.Upsert(
		context.Background(),
		document,
	)

	if err != documents.ErrEmptyDocument {

		t.Fatalf(
			"expected ErrEmptyDocument got %v",
			err,
		)
	}
}