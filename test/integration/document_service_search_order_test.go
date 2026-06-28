package integration

import (
	"context"
	"testing"

	"indexer/internal/documentfilters"
	"indexer/internal/documents"
	"indexer/internal/opensearch"
)

func TestDocumentServiceSearchOrderStringDesc(
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
			"order-string-desc",
		)

		syncDocuments(
			t,
			db,
		)

		_ = db.Close()
	})

	docs := []*documents.Document{
		{
			Namespace:  "order-string-desc",
			ExternalID: "3",
			Payload: []byte(`{
				"name":"Carlos"
			}`),
		},
		{
			Namespace:  "order-string-desc",
			ExternalID: "1",
			Payload: []byte(`{
				"name":"Ana"
			}`),
		},
		{
			Namespace:  "order-string-desc",
			ExternalID: "2",
			Payload: []byte(`{
				"name":"Bruno"
			}`),
		},
	}

	for _, doc := range docs {

		err := documentService.Upsert(
			context.Background(),
			doc,
		)

		if err != nil {
			t.Fatal(err)
		}
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
			nil,
			"",
			&documents.SearchOrder{
				Field:     "name",
				Direction: documents.SortDirectionDesc,
				ValueType: documents.OrderValueTypeString,
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 3 {

		t.Fatalf(
			"expected 3 documents got %d",
			len(results),
		)
	}

	expected := []string{
		"3", // Carlos
		"2", // Bruno
		"1", // Ana
	}

	for i, result := range results {

		if result.ExternalID != expected[i] {

			t.Fatalf(
				"expected external id %s got %s",
				expected[i],
				result.ExternalID,
			)
		}
	}
}