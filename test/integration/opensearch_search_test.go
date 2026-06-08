package integration

import (
	"context"
	"testing"

	"indexer/internal/opensearch"
)

func TestSearchDocument(
	t *testing.T,
) {

	ensureDocumentsIndex(t)

	cfg := createTestConfig()

	client :=
		opensearch.NewClient(cfg)

	document := &opensearch.Document{
		Namespace: "test",
		ExternalID: "999",
		Title: "Magento 2.4",
		Text: "Magento é uma plataforma de ecommerce",
	}

	t.Cleanup(func() {

		_ = opensearch.DeleteDocument(
			context.Background(),
			client,
			"test",
			"999",
		)
	})

	err := opensearch.IndexDocument(
		context.Background(),
		client,
		document,
	)

	if err != nil {
		t.Fatal(err)
	}

	results, err :=
		opensearch.Search(
			context.Background(),
			client,
			"Magento",
			10,
		)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) == 0 {

		t.Fatal(
			"no search results returned",
		)
	}

	found := false

	for _, result := range results {

		if result.ExternalID == "999" {

			found = true

			break
		}
	}

	if !found {

		t.Fatal(
			"indexed document not found",
		)
	}
}