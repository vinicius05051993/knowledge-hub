package integration

import (
	"context"
	"testing"

	"indexer/internal/opensearch"
)

func ensureDocumentsIndex(
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