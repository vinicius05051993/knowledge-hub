package integration

import (
	"context"
	"testing"

	"indexer/internal/documentfilters"
	"indexer/internal/documents"
	"indexer/internal/opensearch"
)

func TestDocumentServiceSearchOrderDateAsc(
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
			"order-date-asc",
		)

		syncDocuments(
			t,
			db,
		)

		_ = db.Close()
	})

	docs := []*documents.Document{
		{
			Namespace:  "order-date-asc",
			ExternalID: "3",
			Payload: []byte(`{
				"created_at":"2025-06-10"
			}`),
		},
		{
			Namespace:  "order-date-asc",
			ExternalID: "1",
			Payload: []byte(`{
				"created_at":"2023-01-15"
			}`),
		},
		{
			Namespace:  "order-date-asc",
			ExternalID: "2",
			Payload: []byte(`{
				"created_at":"2024-03-20"
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
				Field:     "created_at",
				Direction: documents.SortDirectionAsc,
				ValueType: documents.OrderValueTypeDate,
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
		"1",
		"2",
		"3",
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

func TestDocumentServiceSearchOrderDateDesc(
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
			"order-date-desc",
		)

		syncDocuments(
			t,
			db,
		)

		_ = db.Close()
	})

	docs := []*documents.Document{
		{
			Namespace:  "order-date-desc",
			ExternalID: "3",
			Payload: []byte(`{
				"created_at":"2025-06-10"
			}`),
		},
		{
			Namespace:  "order-date-desc",
			ExternalID: "1",
			Payload: []byte(`{
				"created_at":"2023-01-15"
			}`),
		},
		{
			Namespace:  "order-date-desc",
			ExternalID: "2",
			Payload: []byte(`{
				"created_at":"2024-03-20"
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
				Field:     "created_at",
				Direction: documents.SortDirectionDesc,
				ValueType: documents.OrderValueTypeDate,
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
		"3",
		"2",
		"1",
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