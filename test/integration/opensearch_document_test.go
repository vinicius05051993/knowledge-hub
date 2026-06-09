package integration

import (
	"context"
	"testing"

	"indexer/internal/opensearch"
)

func TestIndexDocument(
	t *testing.T,
) {

	ensureDocumentsIndex(t)

	cfg := createTestConfig()

	client :=
		opensearch.NewClient(cfg)

	document := &opensearch.Document{
		DocumentKey: "test:123",
		Namespace: "test",
		ExternalID: "123",
		Title: "Magento 2.4",
		Text: "Magento é uma plataforma de ecommerce",
	}

	t.Cleanup(func() {

		_ = opensearch.DeleteDocument(
			context.Background(),
			client,
			"test",
			"123",
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
}