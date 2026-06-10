package documentfilters

import (
	"context"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(
	db *sqlx.DB,
) *Repository {

	return &Repository{
		db: db,
	}
}

func (r *Repository) Replace(
	ctx context.Context,
	documentKey string,
	filters map[string]string,
) error {

	_, err := r.db.ExecContext(
		ctx,
		`
		DELETE
		FROM document_filters
		WHERE document_key = ?
		`,
		documentKey,
	)

	if err != nil {
		return err
	}

	if len(filters) == 0 {
		return nil
	}

	values := make(
		[]string,
		0,
		len(filters),
	)

	args := make(
		[]any,
		0,
		len(filters)*3,
	)

	for field, value := range filters {

		values = append(
			values,
			"(?,?,?)",
		)

		args = append(
			args,
			documentKey,
			field,
			value,
		)
	}

	query := `
	INSERT INTO document_filters (
		document_key,
		field_name,
		field_value
	)
	VALUES ` + strings.Join(
		values,
		",",
	)

	_, err = r.db.ExecContext(
		ctx,
		query,
		args...,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteByDocumentKey(
	ctx context.Context,
	documentKey string,
) error {

	_, err := r.db.ExecContext(
		ctx,
		`
		DELETE
		FROM document_filters
		WHERE document_key = ?
		`,
		documentKey,
	)

	return err
}