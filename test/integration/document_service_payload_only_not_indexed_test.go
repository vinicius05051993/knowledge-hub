package integration

import (
	"context"
	"testing"

	"indexer/internal/documentfilters"
	"indexer/internal/documents"
	"indexer/internal/opensearch"
)

func TestDocumentServicePayloadOnlyNotIndexed(
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

	namespace := "payload-only-test"

	t.Cleanup(func() {

		_ = documentService.DeleteByNamespace(
			context.Background(),
			namespace,
		)

		syncDocuments(
			t,
			db,
		)

		_ = db.Close()
	})

	err := documentService.Upsert(
		context.Background(),
		&documents.Document{
			Namespace:  namespace,
			ExternalID: "doc-1",
			Payload: []byte(`{
				"sku":"ABC123",
				"brand":"Dell",
				"tag":[
					"tag1",
					"tag2"
				]
			}`),
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	syncDocuments(
		t,
		db,
	)

	document, err :=
		documentRepository.FindByExternalID(
			context.Background(),
			namespace,
			"doc-1",
		)

	if err != nil {
		t.Fatal(err)
	}

	if document.SyncStatus !=
		documents.SyncStatusSynced {

		t.Fatalf(
			"expected Synced got %d",
			document.SyncStatus,
		)
	}

	//
	// Não deve existir no OpenSearch
	//

	results, err :=
		searchService.Search(
			context.Background(),
			"Dell",
			0,
			10,
		)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 0 {

		t.Fatalf(
			"expected 0 indexed documents got %d",
			len(results),
		)
	}

	//
	// Mas deve ser encontrado pelos filtros
	//

	searchResults, err :=
		documentService.Search(
			context.Background(),
			"",
			0,
			10,
			map[string]string{
				"sku": "ABC123",
			},
			documents.FilterTypeAnd,
			nil,
		)

	if err != nil {
		t.Fatal(err)
	}

	if len(searchResults) != 1 {

		t.Fatalf(
			"expected 1 document got %d",
			len(searchResults),
		)
	}

	if searchResults[0].ExternalID !=
		"doc-1" {

		t.Fatalf(
			"unexpected external id %s",
			searchResults[0].ExternalID,
		)
	}
}