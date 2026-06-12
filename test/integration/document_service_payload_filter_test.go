package integration

import (
	"context"
	"testing"

	"indexer/internal/documentfilters"
	"indexer/internal/documents"
	"indexer/internal/opensearch"
)

func TestDocumentServiceSearchPayloadOnlyByFilter(
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
			"payload-filter",
		)

		syncDocuments(
			t,
			db,
		)

		_ = db.Close()
	})

	document := &documents.Document{
		Namespace:  "payload-filter",
		ExternalID: "payload-only",
		Payload: []byte(`{
			"sku":"ABC123",
			"type":"product"
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

	results, err :=
		documentService.Search(
			context.Background(),
			"",
			0,
			10,
			map[string]string{
				"sku": "ABC123",
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {

		t.Fatalf(
			"expected 1 document got %d",
			len(results),
		)
	}

	if results[0].ExternalID !=
		"payload-only" {

		t.Fatalf(
			"unexpected external id %s",
			results[0].ExternalID,
		)
	}
}