package documentfilters

import (
	"context"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

const batchSize = 500

type filterRow struct {
	DocumentKey string
	Field       string
	Value       string
}

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

	return r.ReplaceBatch(
		ctx,
		map[string]map[string]string{
			documentKey: filters,
		},
	)
}

func (r *Repository) DeleteByDocumentKey(
	ctx context.Context,
	documentKey string,
) error {

	query := `
	DELETE
	FROM document_filters
	WHERE document_key = ?
	`

	query = r.db.Rebind(query)

	_, err := r.db.ExecContext(
		ctx,
		query,
		documentKey,
	)

	return err
}

func (r *Repository) DeleteByDocumentKeys(
	ctx context.Context,
	documentKeys []string,
) error {

	if len(documentKeys) == 0 {
		return nil
	}

	query, args, err := sqlx.In(
		`
		DELETE
		FROM document_filters
		WHERE document_key IN (?)
		`,
		documentKeys,
	)

	if err != nil {
		return err
	}

	query = r.db.Rebind(query)

	_, err = r.db.ExecContext(
		ctx,
		query,
		args...,
	)

	return err
}

func (r *Repository) ReplaceBatch(
	ctx context.Context,
	documents map[string]map[string]string,
) error {

	if len(documents) == 0 {
		return nil
	}

	documentKeys := make(
		[]string,
		0,
		len(documents),
	)

	for documentKey := range documents {

		documentKeys = append(
			documentKeys,
			documentKey,
		)
	}

	err := r.DeleteByDocumentKeys(
		ctx,
		documentKeys,
	)

	if err != nil {
		return err
	}

	rows := make(
		[]filterRow,
		0,
	)

	for documentKey, filters := range documents {

		for field, value := range filters {

			rows = append(
				rows,
				filterRow{
					DocumentKey: documentKey,
					Field:       field,
					Value:       value,
				},
			)
		}
	}

	for i := 0; i < len(rows); i += batchSize {

		end := i + batchSize

		if end > len(rows) {
			end = len(rows)
		}

		err := r.replaceChunk(
			ctx,
			rows[i:end],
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) replaceChunk(
	ctx context.Context,
	rows []filterRow,
) error {

	if len(rows) == 0 {
		return nil
	}

	values := make(
		[]string,
		0,
		len(rows),
	)

	args := make(
		[]any,
		0,
		len(rows)*3,
	)

	param := 1

	for _, row := range rows {

		values = append(
			values,
			"(@p"+strconv.Itoa(param)+
				",@p"+strconv.Itoa(param+1)+
				",@p"+strconv.Itoa(param+2)+")",
		)

		args = append(
			args,
			row.DocumentKey,
			row.Field,
			row.Value,
		)

		param += 3
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

	_, err := r.db.ExecContext(
		ctx,
		query,
		args...,
	)

	return err
}