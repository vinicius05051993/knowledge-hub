package documents

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Upsert(
	ctx context.Context,
	doc *Document,
) error {

	query := `
	INSERT INTO documents (
		document_key,
		namespace,
		external_id,
		title,
		text,
		payload,
		created_at,
		updated_at
	)
	VALUES (
		?,?,?,?,?,?,?,?
	)
	ON DUPLICATE KEY UPDATE
		title = VALUES(title),
		text = VALUES(text),
		payload = VALUES(payload),
		updated_at = VALUES(updated_at)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		doc.DocumentKey,
		doc.Namespace,
		doc.ExternalID,
		doc.Title,
		doc.Text,
		doc.Payload,
		doc.CreatedAt,
		doc.UpdatedAt,
	)

	return err
}

func (r *Repository) FindByExternalID(
	ctx context.Context,
	namespace string,
	externalID string,
) (*Document, error) {

	query := `
	SELECT *
	FROM documents
	WHERE namespace = ?
	AND external_id = ?
	LIMIT 1
	`

	var document Document

	err := r.db.GetContext(
		ctx,
		&document,
		query,
		namespace,
		externalID,
	)

	if err != nil {
		return nil, err
	}

	return &document, nil
}

func (r *Repository) FindByDocumentKeys(
	ctx context.Context,
	documentKeys []string,
) ([]Document, error) {

	if len(documentKeys) == 0 {
		return []Document{}, nil
	}

	query, args, err := sqlx.In(
		`
		SELECT *
		FROM documents
		WHERE document_key IN (?)
		`,
		documentKeys,
	)

	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(
		query,
	)

	var documents []Document

	err = r.db.SelectContext(
		ctx,
		&documents,
		query,
		args...,
	)

	if err != nil {
		return nil, err
	}

	return documents, nil
}

func (r *Repository) DeleteByExternalIDs(
	ctx context.Context,
	namespace string,
	externalIDs []string,
) error {

	if len(externalIDs) == 0 {
		return nil
	}

	query, args, err := sqlx.In(
		`
		DELETE FROM documents
		WHERE namespace = ?
		AND external_id IN (?)
		`,
		namespace,
		externalIDs,
	)

	if err != nil {
		return err
	}

	query = r.db.Rebind(
		query,
	)

	_, err = r.db.ExecContext(
		ctx,
		query,
		args...,
	)

	return err
}