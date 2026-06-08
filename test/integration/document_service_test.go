package integration

import (
	"context"
	"testing"

	"indexer/internal/documents"
	"indexer/internal/opensearch"
)

func TestDocumentServiceUpsert(
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

	t.Cleanup(func() {

		_ = opensearch.DeleteDocument(
			context.Background(),
			searchClient,
			"test",
			"service-test",
		)
	})

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

	document := &documents.Document{
		Namespace: "test",
		ExternalID: "service-test",
		Title:      "Magento 2.4",
		Text:       "Magento ecommerce platform",
		Payload: []byte(`{
			"id":123
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
			"test",
			"service-test",
		)

	if err != nil {
		t.Fatal(err)
	}

	if foundDocument == nil {

		t.Fatal(
			"document not found in mysql",
		)
	}

	if foundDocument.Title !=
		"Magento 2.4" {

		t.Fatal(
			"unexpected title",
		)
	}

	results, err :=
		opensearch.Search(
			context.Background(),
			searchClient,
			"Magento",
			10,
		)

	if err != nil {
		t.Fatal(err)
	}

	found := false

	for _, result := range results {

		if result.ExternalID ==
			"service-test" {

			found = true

			break
		}
	}

	if !found {

		t.Fatal(
			"document not found in opensearch",
		)
	}
}