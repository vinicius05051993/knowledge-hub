package documents

import "time"

const (
	SyncStatusSynced = iota
	SyncStatusPendingUpsert
	SyncStatusPendingDelete
	SyncStatusPendingDeindex
)

type Document struct {
	ID int64 `db:"id"`

	DocumentKey string `db:"document_key"`

	Namespace string `db:"namespace"`
	ExternalID string `db:"external_id"`

	Title string `db:"title"`
	Text string `db:"text"`

	Payload []byte `db:"payload"`

	SyncStatus int        `db:"sync_status"`
	DeletedAt  *time.Time `db:"deleted_at"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type SearchDocument struct {
	Document

	Highlights map[string]string
}