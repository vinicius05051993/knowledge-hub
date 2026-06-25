package integration

import (
	"context"
	"encoding/json"
	"testing"

	"indexer/internal/documentfilters"
	"indexer/internal/documents"
	"indexer/internal/opensearch"
)

func TestDocumentPayloadArrayReturned(
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

	namespace := "payload-array-return"

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
				"brand":"Dell",
				"sku":"ABC123",
				"tag":[
					"tag1",
					"tag2",
					"tag3"
				]
			}`),
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	document, err :=
		documentRepository.FindByExternalID(
			context.Background(),
			namespace,
			"doc-1",
		)

	if err != nil {
		t.Fatal(err)
	}

	var payload map[string]any

	err = json.Unmarshal(
		document.Payload,
		&payload,
	)

	if err != nil {
		t.Fatal(err)
	}

	brand, ok :=
		payload["brand"].(string)

	if !ok {

		t.Fatalf(
			"brand should be string",
		)
	}

	if brand != "Dell" {

		t.Fatalf(
			"expected Dell got %s",
			brand,
		)
	}

	sku, ok :=
		payload["sku"].(string)

	if !ok {

		t.Fatalf(
			"sku should be string",
		)
	}

	if sku != "ABC123" {

		t.Fatalf(
			"expected ABC123 got %s",
			sku,
		)
	}

	tags, ok :=
		payload["tag"].([]any)

	if !ok {

		t.Fatalf(
			"tag should be array",
		)
	}

	if len(tags) != 3 {

		t.Fatalf(
			"expected 3 tags got %d",
			len(tags),
		)
	}

	expected := map[string]bool{
		"tag1": false,
		"tag2": false,
		"tag3": false,
	}

	for _, tag := range tags {

		value, ok := tag.(string)

		if !ok {
			t.Fatalf(
				"tag value should be string",
			)
		}

		if _, exists := expected[value]; exists {
			expected[value] = true
		}
	}

	for tag, found := range expected {

		if !found {

			t.Fatalf(
				"missing tag %s",
				tag,
			)
		}
	}
}