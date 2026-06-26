package integration

import (
	"context"
	"testing"

	"indexer/internal/documents"
	"indexer/internal/opensearch"
	"indexer/internal/documentfilters"
)

func TestDocumentServiceSearchWithFilters(
	t *testing.T,
) {

	ensureDocumentsIndex(t)

	db := createDB(t)

	t.Cleanup(func() {

		_ = db.Close()
	})
	
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

	document1 := &documents.Document{
		Namespace:  "test",
		ExternalID: "1",
		Title:      "Magento 2.4",
		Text:       "Magento ecommerce platform",
		Payload: []byte(`{
			"sku":"123",
			"brand":"Dell"
		}`),
	}

	document2 := &documents.Document{
		Namespace:  "test",
		ExternalID: "2",
		Title:      "Magento Commerce",
		Text:       "Adobe Commerce platform",
		Payload: []byte(`{
			"sku":"456",
			"brand":"HP"
		}`),
	}

	err := documentService.Upsert(
		context.Background(),
		document1,
	)

	if err != nil {
		t.Fatal(err)
	}

	err = documentService.Upsert(
		context.Background(),
		document2,
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
			"Magento",
			0,
			10,
			map[string]string{
				"sku": "123",
			},
			"",
			nil,
		)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {

		t.Fatalf(
			"expected 1 document, got %d",
			len(results),
		)
	}

	if results[0].ExternalID !=
		"1" {

		t.Fatal(
			"unexpected document returned",
		)
	}

	if results[0].DocumentKey !=
		"test:1" {

		t.Fatal(
			"unexpected document key",
		)
	}
}

func TestDocumentServiceSearchWithOrFilters(
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

		_ = opensearch.DeleteDocument(
			context.Background(),
			searchClient,
			"test",
			"1",
		)

		_ = opensearch.DeleteDocument(
			context.Background(),
			searchClient,
			"test",
			"2",
		)

		_ = filterRepository.DeleteByDocumentKeys(
			context.Background(),
			[]string{
				"test:1",
				"test:2",
			},
		)

		_ = documentService.DeleteByNamespace(
			context.Background(),
			"test",
		)

		_ = db.Close()
	})

	err := documentService.Upsert(
		context.Background(),
		&documents.Document{
			Namespace:  "test",
			ExternalID: "1",
			Title:      "Magento",
			Text:       "Magento ecommerce",
			Payload: []byte(`{
				"sku":"123",
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
			Namespace:  "test",
			ExternalID: "2",
			Title:      "Magento",
			Text:       "Magento ecommerce",
			Payload: []byte(`{
				"sku":"456",
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
			"Magento",
			0,
			10,
			map[string]string{
				"sku":   "123",
				"brand": "HP",
			},
			documents.FilterTypeOr,
			nil,
		)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 2 {

		t.Fatalf(
			"expected 2 documents, got %d",
			len(results),
		)
	}
}

func TestDocumentServiceSearchWithAndFilters(
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

		_ = opensearch.DeleteDocument(
			context.Background(),
			searchClient,
			"test",
			"1",
		)

		_ = opensearch.DeleteDocument(
			context.Background(),
			searchClient,
			"test",
			"2",
		)

		_ = filterRepository.DeleteByDocumentKeys(
			context.Background(),
			[]string{
				"test:1",
				"test:2",
			},
		)

		_ = documentService.DeleteByNamespace(
			context.Background(),
			"test",
		)

		_ = db.Close()
	})

	err := documentService.Upsert(
		context.Background(),
		&documents.Document{
			Namespace:  "test",
			ExternalID: "1",
			Title:      "Magento",
			Text:       "Magento ecommerce",
			Payload: []byte(`{
				"sku":"123",
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
			Namespace:  "test",
			ExternalID: "2",
			Title:      "Magento",
			Text:       "Magento ecommerce",
			Payload: []byte(`{
				"sku":"456",
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
			"Magento",
			0,
			10,
			map[string]string{
				"sku":   "123",
				"brand": "HP",
			},
			documents.FilterTypeAnd,
			nil,
		)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 0 {

		t.Fatalf(
			"expected 0 documents, got %d",
			len(results),
		)
	}
}

func TestDocumentServiceSearchDefaultFilterType(
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

		_ = opensearch.DeleteDocument(
			context.Background(),
			searchClient,
			"test",
			"1",
		)

		_ = opensearch.DeleteDocument(
			context.Background(),
			searchClient,
			"test",
			"2",
		)

		_ = filterRepository.DeleteByDocumentKeys(
			context.Background(),
			[]string{
				"test:1",
				"test:2",
			},
		)

		_ = documentService.DeleteByNamespace(
			context.Background(),
			"test",
		)

		_ = db.Close()
	})

	err := documentService.Upsert(
		context.Background(),
		&documents.Document{
			Namespace:  "test",
			ExternalID: "1",
			Title:      "Magento",
			Text:       "Magento ecommerce",
			Payload: []byte(`{
				"sku":"123",
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
			Namespace:  "test",
			ExternalID: "2",
			Title:      "Magento",
			Text:       "Magento ecommerce",
			Payload: []byte(`{
				"sku":"456",
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
			"Magento",
			0,
			10,
			map[string]string{
				"sku":   "123",
				"brand": "HP",
			},
			"",
			nil,
		)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 0 {

		t.Fatalf(
			"expected 0 documents, got %d",
			len(results),
		)
	}
}