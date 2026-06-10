package documentfilters

import (
	"context"

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

	for field, value := range filters {

		_, err = r.db.ExecContext(
			ctx,
			`
			INSERT INTO document_filters (
				document_key,
				field_name,
				field_value
			)
			VALUES (
				?,?,?
			)
			`,
			documentKey,
			field,
			value,
		)

		if err != nil {
			return err
		}
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