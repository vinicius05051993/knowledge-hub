package integration

import (
	"context"
	"testing"

	"indexer/internal/documentfilters"
	"indexer/internal/documents"
	"indexer/internal/opensearch"
)

func TestDocumentServiceSearchOrMultipleFields(
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

	namespace := "or-multiple-fields"

	t.Cleanup(func() {

		_ = documentService.DeleteByNamespace(
			context.Background(),
			namespace,
		)

		syncDocuments(
			t,
			db,
		)

		_ = db.Close()
	})

	err := documentService.Upsert(
		context.Background(),
		&documents.Document{
			Namespace:  namespace,
			ExternalID: "doc-tag",
			Payload: []byte(`{
				"tag":[
					"tag1"
				]
			}`),
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	err = documentService.Upsert(
		context.Background(),
		&documents.Document{
			Namespace:  namespace,
			ExternalID: "doc-brand",
			Payload: []byte(`{
				"brand":"Dell"
			}`),
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	err = documentService.Upsert(
		context.Background(),
		&documents.Document{
			Namespace:  namespace,
			ExternalID: "doc-other",
			Payload: []byte(`{
				"brand":"HP"
			}`),
		},
	)

	if err != nil {
		t.Fatal(err)
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
			map[string]string{
				"tag":   "tag1",
				"brand": "Dell",
			},
			documents.FilterTypeOr,
		)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 2 {

		t.Fatalf(
			"expected 2 documents got %d",
			len(results),
		)
	}

	foundTag := false
	foundBrand := false

	for _, result := range results {

		switch result.ExternalID {

		case "doc-tag":
			foundTag = true

		case "doc-brand":
			foundBrand = true
		}
	}

	if !foundTag {

		t.Fatal(
			"doc-tag not found",
		)
	}

	if !foundBrand {

		t.Fatal(
			"doc-brand not found",
		)
	}
}