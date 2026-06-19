package integration

import (
	"context"
	"strings"
	"testing"

	"indexer/internal/chunker"
	"indexer/internal/documents"
	"indexer/internal/documentfilters"
	"indexer/internal/opensearch"
)

func TestDocumentServiceSearchDeduplicatesChunks(
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

		_ = opensearch.DeleteDocuments(
			context.Background(),
			searchClient,
			[]string{
				"test:chunk-doc",
			},
		)

		_ = filterRepository.DeleteByDocumentKeys(
			context.Background(),
			[]string{
				"test:chunk-doc",
			},
		)

		_ = documentService.DeleteByNamespace(
			context.Background(),
			"test",
		)

		_ = db.Close()
	})

	text := strings.Repeat(
		"Chunk ",
		833,
	)

	expectedChunks :=
		len(
			chunker.Split(text),
		)

	if expectedChunks != 5 {

		t.Fatalf(
			"expected test setup to generate 5 chunks got %d",
			expectedChunks,
		)
	}

	document :=
		&documents.Document{
			Namespace:  "test",
			ExternalID: "chunk-doc",
			Title:      "Chunk Test",
			Text:       text,
			Payload: []byte(`{
				"brand":"Chunk"
			}`),
		}

	err :=
		documentService.Upsert(
			context.Background(),
			document,
		)

	if err != nil {
		t.Fatal(err)
	}

	syncDocuments(
		t,
		db,
	)

	searchResults, err :=
		searchService.Search(
			context.Background(),
			"Chunk",
			0,
			50,
		)

	if err != nil {
		t.Fatal(err)
	}

	chunkCount := 0

	for _, result := range searchResults {

		if result.DocumentKey ==
			"test:chunk-doc" {

			chunkCount++
		}
	}

	if chunkCount != 5 {

		t.Fatalf(
			"expected 5 chunks got %d",
			chunkCount,
		)
	}

	results, err :=
		documentService.Search(
			context.Background(),
			"Chunk",
			0,
			20,
			nil,
			"",
		)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {

		t.Fatalf(
			"expected 1 document got %d",
			len(results),
		)
	}

	if results[0].DocumentKey !=
		"test:chunk-doc" {

		t.Fatalf(
			"unexpected document key %s",
			results[0].DocumentKey,
		)
	}
}