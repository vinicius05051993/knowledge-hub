package integration

import (
	"context"
	"testing"

	"indexer/internal/documentfilters"
	"indexer/internal/documents"
	"indexer/internal/opensearch"
)

func TestDocumentServiceSearchAndMissingValue(
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

	namespace := "and-missing-value"

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
				"tag":[
					"tag1"
				]
			}`),
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	err = documentService.Upsert(
		context.Background(),
		&documents.Document{
			Namespace:  namespace,
			ExternalID: "doc-2",
			Payload: []byte(`{
				"tag":[
					"tag5"
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

	results, err :=
		documentService.Search(
			context.Background(),
			"",
			0,
			10,
			map[string]string{
				"tag": "tag1,tag5",
			},
			documents.FilterTypeAnd,
		)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 0 {

		t.Fatalf(
			"expected 0 documents got %d",
			len(results),
		)
	}
}