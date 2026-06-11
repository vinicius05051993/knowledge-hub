package integration

import (
	"context"
	"testing"

	"indexer/internal/documentfilters"
	"indexer/internal/documents"
	"indexer/internal/opensearch"
)

func TestDocumentServiceUpsert(
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

		_ = opensearch.DeleteDocument(
			context.Background(),
			searchClient,
			"test",
			"service-test",
		)

		_ = documentService.DeleteByNamespace(
			context.Background(),
			"test",
		)

		_ = db.Close()
	})

	document := &documents.Document{
		Namespace:  "test",
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

	syncDocuments(
		t,
		db,
	)

	if document.DocumentKey !=
		"test:service-test" {

		t.Fatal(
			"document key not generated",
		)
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

	if foundDocument.DocumentKey !=
		"test:service-test" {

		t.Fatal(
			"unexpected document key",
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
			0,
			10,
		)

	if err != nil {
		t.Fatal(err)
	}

	found := false

	for _, result := range results {

		if result.DocumentKey ==
			"test:service-test" {

			if result.Namespace !=
				"test" {

				t.Fatal(
					"unexpected namespace",
				)
			}

			if result.ExternalID !=
				"service-test" {

				t.Fatal(
					"unexpected external id",
				)
			}

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