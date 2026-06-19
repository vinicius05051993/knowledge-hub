package integration

import (
	"context"
	"testing"

	"indexer/internal/opensearch"
)

func TestBulkIndexDocuments(
	t *testing.T,
) {

	ensureDocumentsIndex(t)

	cfg := createTestConfig()

	client :=
		opensearch.NewClient(cfg)

	documents :=
		[]*opensearch.Document{
			{
				ID:          "test:bulk-1#0",
				DocumentKey: "test:bulk-1",
				Namespace:   "test",
				ExternalID:  "bulk-1",
				Title:       "Magento",
				Text:        "Magento Ecommerce Chunk 0",
			},
			{
				ID:          "test:bulk-1#1",
				DocumentKey: "test:bulk-1",
				Namespace:   "test",
				ExternalID:  "bulk-1",
				Title:       "Magento",
				Text:        "Magento Ecommerce Chunk 1",
			},
			{
				ID:          "test:bulk-2#0",
				DocumentKey: "test:bulk-2",
				Namespace:   "test",
				ExternalID:  "bulk-2",
				Title:       "Magento",
				Text:        "Magento Ecommerce Chunk 0",
			},
			{
				ID:          "test:bulk-2#1",
				DocumentKey: "test:bulk-2",
				Namespace:   "test",
				ExternalID:  "bulk-2",
				Title:       "Magento",
				Text:        "Magento Ecommerce Chunk 1",
			},
		}

	t.Cleanup(func() {

		_ = opensearch.DeleteDocuments(
			context.Background(),
			client,
			[]string{
				"test:bulk-1",
				"test:bulk-2",
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
			"Magento",
			0,
			20,
		)

	if err != nil {
		t.Fatal(err)
	}

	found := map[string]bool{
		"test:bulk-1": false,
		"test:bulk-2": false,
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
}