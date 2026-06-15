package documents

import (
	"context"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
)

const batchSize = 500
const baseSelect = `
SELECT
	d.*,
	COALESCE(
		(
			SELECT JSON_OBJECTAGG(
				df.field_name,
				df.field_value
			)
			FROM document_filters df
			WHERE df.document_key =
				d.document_key
		),
		JSON_OBJECT()
	) AS payload
FROM documents d
`

var validFilterField =
	regexp.MustCompile(
		`^[a-zA-Z0-9_]+$`,
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

    return r.UpsertBatch(
        ctx,
        []*Document{doc},
    )
}

func (r *Repository) FindByExternalID(
	ctx context.Context,
	namespace string,
	externalID string,
) (*Document, error) {

	query := baseSelect + `
	WHERE d.namespace = ?
	AND d.external_id = ?
	AND d.deleted_at IS NULL
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
		baseSelect+`
		WHERE d.document_key IN (?)
		AND d.deleted_at IS NULL
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

func (r *Repository) DeleteByNamespace(
	ctx context.Context,
	namespace string,
) error {

	_, err := r.db.ExecContext(
		ctx,
		`
		UPDATE documents
		SET
			sync_status = ?,
			deleted_at = NOW()
		WHERE namespace = ?
		`,
		SyncStatusPendingDelete,
		namespace,
	)

	return err
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
		UPDATE documents
		SET
			sync_status = ?,
			deleted_at = NOW()
		WHERE namespace = ?
		AND external_id IN (?)
		`,
		SyncStatusPendingDelete,
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

func (r *Repository) Search(
	ctx context.Context,
	documentKeys []string,
	offset int,
	limit int,
	filters map[string]string,
	filterType string,
) ([]Document, error) {

	query := baseSelect + `
	WHERE 1 = 1
	AND d.deleted_at IS NULL
	`

	args := make(
		[]any,
		0,
	)

	if len(documentKeys) > 0 {

		inQuery, inArgs, err := sqlx.In(
			`
			d.document_key IN (?)
			`,
			documentKeys,
		)

		if err != nil {
			return nil, err
		}

		query += `
		AND ` + inQuery

		args = append(
			args,
			inArgs...,
		)
	}

	if filterType != FilterTypeOr {
		filterType = FilterTypeAnd
	}

	if len(filters) > 0 {

		if filterType == FilterTypeAnd {

			for field, value := range filters {

				if !validFilterField.MatchString(
					field,
				) {
					continue
				}

				query += `
				AND EXISTS (
					SELECT 1
					FROM document_filters f
					WHERE
						f.document_key =
							d.document_key
					AND f.field_name = ?
					AND f.field_value = ?
				)
				`

				args = append(
					args,
					field,
					value,
				)
			}

		} else {

			query += `
			AND (
			`

			first := true

			for field, value := range filters {

				if !validFilterField.MatchString(
					field,
				) {
					continue
				}

				if !first {

					query += `
					OR
					`
				}

				query += `
				EXISTS (
					SELECT 1
					FROM document_filters f
					WHERE
						f.document_key =
							d.document_key
					AND f.field_name = ?
					AND f.field_value = ?
				)
				`

				args = append(
					args,
					field,
					value,
				)

				first = false
			}

			query += `
			)
			`
		}
	}

	if len(documentKeys) == 0 {

		query += `
		LIMIT ?
		OFFSET ?
		`

		args = append(
			args,
			limit,
			offset,
		)
	}

	query = r.db.Rebind(
		query,
	)

	var documents []Document

	err := r.db.SelectContext(
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

func (r *Repository) FindPendingUpserts(
	ctx context.Context,
	limit int,
) ([]Document, error) {

	query := `
	SELECT *
	FROM documents
	WHERE sync_status = ?
	ORDER BY id
	LIMIT ?
	`

	var documents []Document

	err := r.db.SelectContext(
		ctx,
		&documents,
		query,
		SyncStatusPendingUpsert,
		limit,
	)

	if err != nil {
		return nil, err
	}

	return documents, nil
}

func (r *Repository) FindPendingDeletes(
	ctx context.Context,
	limit int,
) ([]Document, error) {

	query := `
	SELECT *
	FROM documents
	WHERE sync_status = ?
	ORDER BY id
	LIMIT ?
	`

	var documents []Document

	err := r.db.SelectContext(
		ctx,
		&documents,
		query,
		SyncStatusPendingDelete,
		limit,
	)

	if err != nil {
		return nil, err
	}

	return documents, nil
}

func (r *Repository) CountPendingUpserts(
	ctx context.Context,
) (int, error) {

	var count int

	err := r.db.GetContext(
		ctx,
		&count,
		`
		SELECT COUNT(*)
		FROM documents
		WHERE sync_status = ?
		`,
		SyncStatusPendingUpsert,
	)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Repository) CountPendingDeletes(
	ctx context.Context,
) (int, error) {

	var count int

	err := r.db.GetContext(
		ctx,
		&count,
		`
		SELECT COUNT(*)
		FROM documents
		WHERE sync_status = ?
		`,
		SyncStatusPendingDelete,
	)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Repository) MarkSyncedByDocumentKeys(
	ctx context.Context,
	documentKeys []string,
) error {

	if len(documentKeys) == 0 {
		return nil
	}

	query, args, err := sqlx.In(
		`
		UPDATE documents
		SET
			sync_status = ?,
			deleted_at = NULL
		WHERE document_key IN (?)
		`,
		SyncStatusSynced,
		documentKeys,
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
		FROM documents
		WHERE document_key IN (?)
		`,
		documentKeys,
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

func (r *Repository) UpsertBatch(
    ctx context.Context,
    docs []*Document,
) error {

    if len(docs) == 0 {
        return nil
    }

    for i := 0; i < len(docs); i += batchSize {

        end := i + batchSize

        if end > len(docs) {
            end = len(docs)
        }

        if err := r.upsertChunk(
            ctx,
            docs[i:end],
        ); err != nil {
            return err
        }
    }

    return nil
}

func (r *Repository) upsertChunk(
	ctx context.Context,
	docs []*Document,
) error {

	values := make(
		[]string,
		0,
		len(docs),
	)

	args := make(
		[]any,
		0,
		len(docs)*9,
	)

	for _, doc := range docs {

		values = append(
			values,
			"(?,?,?,?,?,?,?,?,?)",
		)

		args = append(
			args,
			doc.DocumentKey,
			doc.Namespace,
			doc.ExternalID,
			doc.Title,
			doc.Text,
			doc.SyncStatus,
			doc.DeletedAt,
			doc.CreatedAt,
			doc.UpdatedAt,
		)
	}

	query := `
	INSERT INTO documents (
		document_key,
		namespace,
		external_id,
		title,
		text,
		sync_status,
		deleted_at,
		created_at,
		updated_at
	)
	VALUES ` + strings.Join(values, ",") + `
	ON DUPLICATE KEY UPDATE
		title = VALUES(title),
		text = VALUES(text),
		sync_status = VALUES(sync_status),
		deleted_at = VALUES(deleted_at),
		updated_at = VALUES(updated_at)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		args...,
	)

	return err
}