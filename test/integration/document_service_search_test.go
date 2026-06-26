package integration

import (
	"context"
	"testing"

	"indexer/internal/documents"
	"indexer/internal/opensearch"
	"indexer/internal/documentfilters"
)

func TestDocumentServiceSearch(
	t *testing.T,
) {

	ensureDocumentsIndex(t)

	db := createDB(t)

	cfg := createTestConfig()

	searchClient :=
		opensearch.NewClient(cfg)

	t.Cleanup(func() {

		_ = opensearch.DeleteDocument(
			context.Background(),
			searchClient,
			"test",
			"1",
		)

		_ = opensearch.DeleteDocument(
			context.Background(),
			searchClient,
			"test",
			"2",
		)

		_ = db.Close()
	})

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

	document1 := &documents.Document{
		Namespace:  "test",
		ExternalID: "1",
		Title:      "Adobe 2.4",
		Text:       "Adobe ecommerce platform",
		Payload:    []byte(`{"id":1}`),
	}

	document2 := &documents.Document{
		Namespace:  "test",
		ExternalID: "2",
		Title:      "Adobe Commerce",
		Text:       "Adobe Commerce platform",
		Payload:    []byte(`{"id":2}`),
	}

	err := documentService.Upsert(
		context.Background(),
		document1,
	)

	if err != nil {
		t.Fatal(err)
	}

	err = documentService.Upsert(
		context.Background(),
		document2,
	)

	if err != nil {
		t.Fatal(err)
	}

	syncDocuments(
		t,
		db,
	)

	results, err :=
		documentService.Search(
			context.Background(),
			"Adobe",
			0,
			10,
			nil,
			"",
			nil,
		)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 2 {

		t.Fatalf(
			"expected 2 documents, got %d",
			len(results),
		)
	}
}