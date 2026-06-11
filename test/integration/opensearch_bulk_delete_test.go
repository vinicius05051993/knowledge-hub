package integration

import (
	"context"
	"testing"

	"indexer/internal/opensearch"
)

func TestBulkDeleteDocuments(
	t *testing.T,
) {

	ensureDocumentsIndex(t)

	cfg := createTestConfig()

	client :=
		opensearch.NewClient(cfg)

	documents :=
		[]*opensearch.Document{
			{
				DocumentKey: "test:delete-1",
				Namespace:   "test",
				ExternalID:  "delete-1",
				Title:       "Delete Test",
				Text:        "Document One",
			},
			{
				DocumentKey: "test:delete-2",
				Namespace:   "test",
				ExternalID:  "delete-2",
				Title:       "Delete Test",
				Text:        "Document Two",
			},
		}

	t.Cleanup(func() {

		_ = opensearch.DeleteDocuments(
			context.Background(),
			client,
			[]string{
				"test:delete-1",
				"test:delete-2",
			},
		)
	})

	err :=
		opensearch.BulkIndexDocuments(
			context.Background(),
			client,
			documents,
		)

	if err != nil {
		t.Fatal(err)
	}

	results, err :=
		opensearch.Search(
			context.Background(),
			client,
			"Delete",
			0,
			10,
		)

	if err != nil {
		t.Fatal(err)
	}

	found := 0

	for _, result := range results {

		if result.DocumentKey ==
			"test:delete-1" ||
			result.DocumentKey ==
				"test:delete-2" {

			found++
		}
	}

	if found != 2 {

		t.Fatalf(
			"expected 2 indexed documents got %d",
			found,
		)
	}

	err = opensearch.DeleteDocuments(
		context.Background(),
		client,
		[]string{
			"test:delete-1",
			"test:delete-2",
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	results, err =
		opensearch.Search(
			context.Background(),
			client,
			"Delete",
			0,
			10,
		)

	if err != nil {
		t.Fatal(err)
	}

	for _, result := range results {

		if result.DocumentKey ==
			"test:delete-1" ||
			result.DocumentKey ==
				"test:delete-2" {

			t.Fatalf(
				"document %s still exists",
				result.DocumentKey,
			)
		}
	}
}