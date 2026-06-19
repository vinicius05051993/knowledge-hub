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
				ID:          "test:delete-1#0",
				DocumentKey: "test:delete-1",
				Namespace:   "test",
				ExternalID:  "delete-1",
				Title:       "Delete Test",
				Text:        "Document One Chunk 0",
			},
			{
				ID:          "test:delete-1#1",
				DocumentKey: "test:delete-1",
				Namespace:   "test",
				ExternalID:  "delete-1",
				Title:       "Delete Test",
				Text:        "Document One Chunk 1",
			},
			{
				ID:          "test:delete-2#0",
				DocumentKey: "test:delete-2",
				Namespace:   "test",
				ExternalID:  "delete-2",
				Title:       "Delete Test",
				Text:        "Document Two Chunk 0",
			},
			{
				ID:          "test:delete-2#1",
				DocumentKey: "test:delete-2",
				Namespace:   "test",
				ExternalID:  "delete-2",
				Title:       "Delete Test",
				Text:        "Document Two Chunk 1",
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
			20,
		)

	if err != nil {
		t.Fatal(err)
	}

	found := map[string]bool{
		"test:delete-1": false,
		"test:delete-2": false,
	}

	for _, result := range results {

		if _, ok :=
			found[result.DocumentKey]; ok {

			found[result.DocumentKey] = true
		}
	}

	for documentKey, exists := range found {

		if !exists {

			t.Fatalf(
				"document %s not found",
				documentKey,
			)
		}
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
			20,
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