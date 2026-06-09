package integration

import (
	"context"
	"testing"

	"indexer/internal/documents"
	"indexer/internal/opensearch"
)

func TestDocumentServiceSearch(
	t *testing.T,
) {

	ensureDocumentsIndex(t)

	db := createDB(t)

	defer db.Close()

	cleanupTestData(
		t,
		db,
	)

	cfg := createTestConfig()

	searchClient :=
		opensearch.NewClient(cfg)

	searchService :=
		opensearch.NewService(
			searchClient,
		)

	documentRepository :=
		documents.NewRepository(db)

	documentService :=
		documents.NewService(
			documentRepository,
			searchService,
		)

	document1 := &documents.Document{
		Namespace:  "test",
		ExternalID: "1",
		Title:      "Magento 2.4",
		Text:       "Magento ecommerce platform",
		Payload:    []byte(`{"id":1}`),
	}

	document2 := &documents.Document{
		Namespace:  "test",
		ExternalID: "2",
		Title:      "Magento Commerce",
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

	results, err :=
		documentService.Search(
			context.Background(),
			"Magento",
			10,
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