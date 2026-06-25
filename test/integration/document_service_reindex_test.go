package integration

import (
	"context"
	"testing"

	"indexer/internal/documentfilters"
	"indexer/internal/documents"
	"indexer/internal/opensearch"
)

func TestDocumentServiceReindex(
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

	namespace := "reindex-test"

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
	// Cria documento indexável
	//

	err := documentService.Upsert(
		context.Background(),
		&documents.Document{
			Namespace:  namespace,
			ExternalID: "doc-1",
			Title:      "Notebook Dell",
			Text:       "Core i7",
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	syncDocuments(
		t,
		db,
	)

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
	// Remove title/text
	//

	err = documentService.Upsert(
		context.Background(),
		&documents.Document{
			Namespace:  namespace,
			ExternalID: "doc-1",
			Payload: []byte(`{
				"sku":"ABC123"
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
	// Volta a possuir conteúdo pesquisável
	//

	err = documentService.Upsert(
		context.Background(),
		&documents.Document{
			Namespace:  namespace,
			ExternalID: "doc-1",
			Title:      "Notebook Lenovo",
			Text:       "Ryzen 7",
			Payload: []byte(`{
				"sku":"ABC123"
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
		documents.SyncStatusPendingUpsert {

		t.Fatalf(
			"expected PendingUpsert got %d",
			document.SyncStatus,
		)
	}

	syncDocuments(
		t,
		db,
	)

	results, err =
		searchService.Search(
			context.Background(),
			"Lenovo",
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

	if results[0].DocumentKey !=
		namespace+":doc-1" {

		t.Fatalf(
			"unexpected document key %s",
			results[0].DocumentKey,
		)
	}

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
}