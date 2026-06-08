package apikeys

import "time"

type APIKey struct {
	ID int64 `db:"id"`

	Name string `db:"name"`

	Namespace string `db:"namespace"`

	APIKeyHash string `db:"api_key_hash"`

	Permissions string `db:"permissions"`

	Active bool `db:"active"`

	CreatedAt time.Time `db:"created_at"`
}