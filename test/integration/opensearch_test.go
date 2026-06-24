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
		)

	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateDocumentsIndexTwice(
	t *testing.T,
) {

	cfg := createTestConfig()

	client :=
		opensearch.NewClient(cfg)

	err :=
		opensearch.CreateDocumentsIndex(
			context.Background(),
			client,
		)

	if err != nil {
		t.Fatal(err)
	}

	err =
		opensearch.CreateDocumentsIndex(
			context.Background(),
			client,
		)

	if err != nil {
		t.Fatal(err)
	}
}