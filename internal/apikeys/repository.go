package apikeys

import (
	"context"
	"fmt"

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

func (r *Repository) FindByHash(
	ctx context.Context,
	hash string,
) (*APIKey, error) {

	query := `
	SELECT *
	FROM api_keys
	WHERE api_key_hash = ?
	AND active = TRUE
	LIMIT 1
	`

	var apiKey APIKey

	err := r.db.GetContext(
		ctx,
		&apiKey,
		query,
		hash,
	)

	if err != nil {
		return nil, err
	}

	return &apiKey, nil
}

func (r *Repository) Create(
	ctx context.Context,
	apiKey *APIKey,
) error {

	query := `
	INSERT INTO api_keys (
		name,
		namespace,
		api_key_hash,
		permissions,
		active,
		created_at
	)
	VALUES (
		?,?,?,?,?,?
	)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		apiKey.Name,
		apiKey.Namespace,
		apiKey.APIKeyHash,
		apiKey.Permissions,
		apiKey.Active,
		apiKey.CreatedAt,
	)

	return err
}

func (r *Repository) DeleteByNamespace(
	ctx context.Context,
	namespace string,
) error {

	_, err := r.db.ExecContext(
		ctx,
		`
		DELETE
		FROM api_keys
		WHERE namespace = ?
		`,
		namespace,
	)

	if err != nil {

		fmt.Println(
			"delete error:",
			err,
		)

		return err
	}

	return nil
}