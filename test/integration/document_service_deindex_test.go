package integration

import (
	"context"
	"testing"

	"indexer/internal/documentfilters"
	"indexer/internal/documents"
	"indexer/internal/opensearch"
)

func TestDocumentServicePendingDeindex(
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

	namespace := "pending-deindex-test"

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

	//
	// Cria um documento indexável
	//

	err := documentService.Upsert(
		context.Background(),
		&documents.Document{
			Namespace:  namespace,
			ExternalID: "doc-1",
			Title:      "Notebook Dell",
			Text:       "Notebook i7 16GB",
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	syncDocuments(
		t,
		db,
	)

	//
	// Confirma que está no OpenSearch
	//

	results, err :=
		searchService.Search(
			context.Background(),
			"Notebook",
			0,
			10,
		)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {

		t.Fatalf(
			"expected 1 indexed document got %d",
			len(results),
		)
	}

	//
	// Atualiza removendo title/text
	//

	err = documentService.Upsert(
		context.Background(),
		&documents.Document{
			Namespace:  namespace,
			ExternalID: "doc-1",
			Payload: []byte(`{
				"sku":"ABC123",
				"brand":"Dell"
			}`),
		},
	)

	if err != nil {
		t.Fatal(err)
	}

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
		documents.SyncStatusPendingDeindex {

		t.Fatalf(
			"expected PendingDeindex got %d",
			document.SyncStatus,
		)
	}

	//
	// Processa o PendingDeindex
	//

	syncDocuments(
		t,
		db,
	)

	//
	// Deve ter saído do OpenSearch
	//

	results, err =
		searchService.Search(
			context.Background(),
			"Notebook",
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
	// Mas continua existindo no banco
	//

	document, err =
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
	// Continua pesquisável pelos filtros
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
		)

	if err != nil {
		t.Fatal(err)
	}

	if len(searchResults) != 1 {

		t.Fatalf(
			"expected 1 filtered document got %d",
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