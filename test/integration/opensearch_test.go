package integration

import (
	"context"
	"testing"

	"indexer/internal/opensearch"
)

func TestCreateDocumentsIndex(
	t *testing.T,
) {

	cfg := createTestConfig()

	client :=
		opensearch.NewClient(cfg)

	err :=
		opensearch.CreateDocumentsIndex(
			context.Background(),
			client,
			cfg,
		)

	if err != nil {
		t.Fatal(err)
	}
}

func TestDocumentsIndexExists(
	t *testing.T,
) {

	ensureDocumentsIndex(t)

	cfg := createTestConfig()

	client :=
		opensearch.NewClient(cfg)

	exists, err :=
		opensearch.IndexExists(
			context.Background(),
			client,
			opensearch.DocumentsIndex,
		)

	if err != nil {
		t.Fatal(err)
	}

	if !exists {

		t.Fatal(
			"documents index does not exist",
		)
	}
}