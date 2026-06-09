package integration

import (
	"context"
	"testing"

	"indexer/internal/documents"
	"indexer/internal/opensearch"
)

func TestDocumentServiceDelete(
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

	document := &documents.Document{
		Namespace:  "test",
		ExternalID: "delete-test",
		Title:      "Magento Delete",
		Text:       "Document to be deleted",
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

	err = documentService.Delete(
		context.Background(),
		"test",
		[]string{
			"delete-test",
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	_, err =
		documentRepository.FindByExternalID(
			context.Background(),
			"test",
			"delete-test",
		)

	if err == nil {

		t.Fatal(
			"document should not exist in mysql",
		)
	}

	results, err :=
		opensearch.Search(
			context.Background(),
			searchClient,
			"Magento Delete",
			10,
		)

	if err != nil {
		t.Fatal(err)
	}

	for _, result := range results {

		if result.DocumentKey ==
			"test:delete-test" {

			t.Fatal(
				"document should not exist in opensearch",
			)
		}
	}
}