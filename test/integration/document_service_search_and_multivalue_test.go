package integration

import (
	"context"
	"testing"

	"indexer/internal/documentfilters"
	"indexer/internal/documents"
	"indexer/internal/opensearch"
)

func TestDocumentServiceSearchAndMultiValue(
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

	namespace := "and-multi-value"

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
			ExternalID: "doc-1",
			Payload: []byte(`{
				"tag":[
					"tag1",
					"tag5"
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
			ExternalID: "doc-2",
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
			ExternalID: "doc-3",
			Payload: []byte(`{
				"tag":[
					"tag5"
				]
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
			[]documents.SearchFilter{
				{
					Field:     "tag",
					Operator:  documents.FilterOperatorEqual,
					Value:     "tag1",
					ValueType: documents.OrderValueTypeString,
				},
				{
					Field:     "tag",
					Operator:  documents.FilterOperatorEqual,
					Value:     "tag5",
					ValueType: documents.OrderValueTypeString,
				},
			},
			documents.FilterTypeAnd,
			nil,
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

	if results[0].ExternalID !=
		"doc-1" {

		t.Fatalf(
			"expected doc-1 got %s",
			results[0].ExternalID,
		)
	}
}