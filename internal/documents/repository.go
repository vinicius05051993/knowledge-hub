package documents

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/jmoiron/sqlx"
)

const batchSize = 500
const baseSelect = `
SELECT
	d.*,
	COALESCE(
	(
		SELECT
			'{' +
			ISNULL(
				STRING_AGG(
					field_json,
					','
				),
				''
			) +
			'}'
		FROM (
			SELECT
				'"' +
				STRING_ESCAPE(
					field_name,
					'json'
				) +
				'":' +
				CASE
					WHEN COUNT(*) = 1 THEN
						'"' +
						STRING_ESCAPE(
							MAX(field_value),
							'json'
						) +
						'"'
					ELSE
						'[' +
						STRING_AGG(
							'"' +
							STRING_ESCAPE(
								field_value,
								'json'
							) +
							'"',
							','
						) +
						']'
				END AS field_json
			FROM document_filters
			WHERE document_key = d.document_key
			GROUP BY field_name
		) grouped_fields
	),
	'{}'
	) AS payload
FROM documents d
`

const baseSelectTop1 = `
SELECT TOP 1
	d.*,
	COALESCE(
	(
		SELECT
			'{' +
			ISNULL(
				STRING_AGG(
					field_json,
					','
				),
				''
			) +
			'}'
		FROM (
			SELECT
				'"' +
				STRING_ESCAPE(
					field_name,
					'json'
				) +
				'":' +
				CASE
					WHEN COUNT(*) = 1 THEN
						'"' +
						STRING_ESCAPE(
							MAX(field_value),
							'json'
						) +
						'"'
					ELSE
						'[' +
						STRING_AGG(
							'"' +
							STRING_ESCAPE(
								field_value,
								'json'
							) +
							'"',
							','
						) +
						']'
				END AS field_json
			FROM document_filters
			WHERE document_key = d.document_key
			GROUP BY field_name
		) grouped_fields
	),
	'{}'
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

	query := baseSelectTop1 + `
	WHERE d.namespace = ?
	AND d.external_id = ?
	AND d.deleted_at IS NULL
	`

	query = r.db.Rebind(query)

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

	query := `
	UPDATE documents
	SET
		sync_status = ?,
		deleted_at = GETDATE()
	WHERE namespace = ?
	`

	query = r.db.Rebind(query)

	_, err := r.db.ExecContext(
		ctx,
		query,
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
			deleted_at = GETDATE()
		WHERE namespace = ?
		AND external_id IN (?)
		`,
		SyncStatusPendingDelete,
		namespace,
		externalIDs,
	)

	if err != nil {

		log.Printf(
			"DeleteByExternalIDs sqlx.In failed\nnamespace=%s\nexternalIDs=%d\nerror=%v",
			namespace,
			len(externalIDs),
			err,
		)

		return fmt.Errorf(
			"documents.Repository.DeleteByExternalIDs.sqlxIn: %w",
			err,
		)
	}

	query = r.db.Rebind(
		query,
	)

	_, err = r.db.ExecContext(
		ctx,
		query,
		args...,
	)

	if err != nil {

		log.Printf(
			"DeleteByExternalIDs failed\nnamespace=%s\nexternalIDs=%d\nquery:\n%s\nargs=%d\nerror=%v",
			namespace,
			len(externalIDs),
			query,
			len(args),
			err,
		)

		return fmt.Errorf(
			"documents.Repository.DeleteByExternalIDs: %w",
			err,
		)
	}

	return nil
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
		ORDER BY d.id
		OFFSET ? ROWS
		FETCH NEXT ? ROWS ONLY
		`

		args = append(
			args,
			offset,
			limit,
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
	SELECT TOP ` + strconv.Itoa(limit) + ` *
	FROM documents
	WHERE sync_status = ?
	ORDER BY id
	`

	query = r.db.Rebind(query)

	var documents []Document

	err := r.db.SelectContext(
		ctx,
		&documents,
		query,
		SyncStatusPendingUpsert,
	)

	if err != nil {

		log.Printf(
			"FindPendingUpserts failed\nlimit=%d\nquery:\n%s\nerror:%v",
			limit,
			query,
			err,
		)

		return nil, fmt.Errorf(
			"documents.Repository.FindPendingUpserts: %w",
			err,
		)
	}

	return documents, nil
}

func (r *Repository) FindPendingDeletes(
	ctx context.Context,
	limit int,
) ([]Document, error) {

	query := `
	SELECT TOP ` + strconv.Itoa(limit) + ` *
	FROM documents
	WHERE sync_status = ?
	ORDER BY id
	`

	query = r.db.Rebind(query)

	var documents []Document

	err := r.db.SelectContext(
		ctx,
		&documents,
		query,
		SyncStatusPendingDelete,
	)

	if err != nil {

		log.Printf(
			"FindPendingDeletes failed\nlimit=%d\nquery:\n%s\nerror:%v",
			limit,
			query,
			err,
		)

		return nil, fmt.Errorf(
			"documents.Repository.FindPendingDeletes: %w",
			err,
		)
	}

	return documents, nil
}

func (r *Repository) CountPendingUpserts(
	ctx context.Context,
) (int, error) {

	var count int

	query := `
	SELECT COUNT(*)
	FROM documents
	WHERE sync_status = ?
	`

	query = r.db.Rebind(query)

	err := r.db.GetContext(
		ctx,
		&count,
		query,
		SyncStatusPendingUpsert,
	)

	if err != nil {

		log.Printf(
			"CountPendingUpserts failed\nquery:\n%s\nerror:%v",
			query,
			err,
		)

		return 0, fmt.Errorf(
			"documents.Repository.CountPendingUpserts: %w",
			err,
		)
	}

	return count, nil
}

func (r *Repository) CountPendingDeletes(
	ctx context.Context,
) (int, error) {

	var count int

	query := `
	SELECT COUNT(*)
	FROM documents
	WHERE sync_status = ?
	`

	query = r.db.Rebind(query)

	err := r.db.GetContext(
		ctx,
		&count,
		query,
		SyncStatusPendingDelete,
	)

	if err != nil {

		log.Printf(
			"CountPendingDeletes failed\nquery:\n%s\nerror:%v",
			query,
			err,
		)

		return 0, fmt.Errorf(
			"documents.Repository.CountPendingDeletes: %w",
			err,
		)
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

		log.Printf(
			"MarkSyncedByDocumentKeys sqlx.In failed\nkeys=%d\nerror=%v",
			len(documentKeys),
			err,
		)

		return fmt.Errorf(
			"documents.Repository.MarkSyncedByDocumentKeys.sqlxIn: %w",
			err,
		)
	}

	query = r.db.Rebind(
		query,
	)

	_, err = r.db.ExecContext(
		ctx,
		query,
		args...,
	)

	if err != nil {

		log.Printf(
			"MarkSyncedByDocumentKeys failed\nkeys=%d\nquery:\n%s\nargs=%d\nerror=%v",
			len(documentKeys),
			query,
			len(args),
			err,
		)

		return fmt.Errorf(
			"documents.Repository.MarkSyncedByDocumentKeys: %w",
			err,
		)
	}

	return nil
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

		log.Printf(
			"DeleteByDocumentKeys sqlx.In failed\nkeys=%d\nerror=%v",
			len(documentKeys),
			err,
		)

		return fmt.Errorf(
			"documents.Repository.DeleteByDocumentKeys.sqlxIn: %w",
			err,
		)
	}

	query = r.db.Rebind(
		query,
	)

	_, err = r.db.ExecContext(
		ctx,
		query,
		args...,
	)

	if err != nil {

		log.Printf(
			"DeleteByDocumentKeys failed\nkeys=%d\nquery:\n%s\nargs=%d\nerror=%v",
			len(documentKeys),
			query,
			len(args),
			err,
		)

		return fmt.Errorf(
			"documents.Repository.DeleteByDocumentKeys: %w",
			err,
		)
	}

	return nil
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

	query := `
	MERGE documents AS target
	USING (
		VALUES (
			@p1, @p2, @p3, @p4, @p5,
			@p6, @p7, @p8, @p9
		)
	) AS source (
		document_key,
		namespace,
		external_id,
		title,
		document_text,
		sync_status,
		deleted_at,
		created_at,
		updated_at
	)
	ON target.document_key = source.document_key

	WHEN MATCHED THEN
		UPDATE SET
			namespace = source.namespace,
			external_id = source.external_id,
			title = source.title,
			[text] = source.document_text,
			sync_status = source.sync_status,
			deleted_at = source.deleted_at,
			updated_at = source.updated_at

	WHEN NOT MATCHED THEN
		INSERT (
			document_key,
			namespace,
			external_id,
			title,
			[text],
			sync_status,
			deleted_at,
			created_at,
			updated_at
		)
		VALUES (
			source.document_key,
			source.namespace,
			source.external_id,
			source.title,
			source.document_text,
			source.sync_status,
			source.deleted_at,
			source.created_at,
			source.updated_at
		);
	`

	for _, doc := range docs {

		_, err := r.db.ExecContext(
			ctx,
			query,
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

		if err != nil {

			log.Printf(
			    "upsertChunk failed | documentKey=%s | namespace=%s | externalID=%s | error=%v",
			    doc.DocumentKey,
			    doc.Namespace,
			    doc.ExternalID,
			    err,
			)

			return fmt.Errorf(
				"documents.Repository.upsertChunk documentKey=%s: %w",
				doc.DocumentKey,
				err,
			)
		}
	}

	return nil
}