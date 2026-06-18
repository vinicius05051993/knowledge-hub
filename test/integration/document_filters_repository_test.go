package integration

import (
	"context"
	"testing"

	"indexer/internal/documentfilters"
	"indexer/internal/documents"
	"indexer/internal/opensearch"
)

func TestDocumentFiltersRepository(
	t *testing.T,
) {

	ensureDocumentsIndex(t)

	db := createDB(t)

	cfg := createTestConfig()

	searchClient :=
		opensearch.NewClient(cfg)

	documentRepository :=
		documents.NewRepository(db)

	filterRepository :=
		documentfilters.NewRepository(
			db,
		)

	searchService :=
		opensearch.NewService(
			searchClient,
		)

	documentService :=
		documents.NewService(
			documentRepository,
			filterRepository,
			searchService,
		)

	namespace := "test"

	t.Cleanup(func() {

		_ = opensearch.DeleteDocument(
			context.Background(),
			searchClient,
			namespace,
			"filters-test",
		)

		_ = filterRepository.DeleteByDocumentKeys(
			context.Background(),
			[]string{
				"test:filters-test",
			},
		)

		_ = documentService.DeleteByNamespace(
			context.Background(),
			namespace,
		)

		_ = db.Close()
	})

	err := documentService.Upsert(
		context.Background(),
		&documents.Document{
			Namespace:  namespace,
			ExternalID: "filters-test",
			Title:      "Produto",
			Text:       "Produto de teste",
			Payload: []byte(`{
				"sku":"vini",
				"brand":"apple",
				"price":100
			}`),
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	type FilterRow struct {
		DocumentKey string `db:"document_key"`
		FieldName   string `db:"field_name"`
		FieldValue  string `db:"field_value"`
	}

	var rows []FilterRow

	query := `
	SELECT
		document_key,
		field_name,
		field_value
	FROM document_filters
	WHERE document_key = ?
	`

	query = db.Rebind(query)

	err = db.Select(
		&rows,
		query,
		"test:filters-test",
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(rows) != 2 {

		t.Fatalf(
			"expected 2 filters got %d",
			len(rows),
		)
	}

	foundSKU := false
	foundBrand := false

	for _, row := range rows {

		if row.FieldName == "sku" &&
			row.FieldValue == "vini" {

			foundSKU = true
		}

		if row.FieldName == "brand" &&
			row.FieldValue == "apple" {

			foundBrand = true
		}
	}

	if !foundSKU {

		t.Fatal(
			"sku filter not found",
		)
	}

	if !foundBrand {

		t.Fatal(
			"brand filter not found",
		)
	}
}