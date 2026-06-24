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
	SELECT TOP 1 *
	FROM api_keys
	WHERE api_key_hash = ?
	AND active = 1
	`

	query = r.db.Rebind(query)

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

	query = r.db.Rebind(query)

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

func (r *Repository) UpdateHashByNamespace(
	ctx context.Context,
	namespace string,
	hash string,
) error {

	query := `
	UPDATE api_keys
	SET api_key_hash = ?
	WHERE namespace = ?
	`

	query = r.db.Rebind(query)

	result, err := r.db.ExecContext(
		ctx,
		query,
		hash,
		namespace,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf(
			"namespace '%s' not found",
			namespace,
		)
	}

	return nil
}

func (r *Repository) DeleteByNamespace(
	ctx context.Context,
	namespace string,
) error {

	query := `
	DELETE
	FROM api_keys
	WHERE namespace = ?
	`

	query = r.db.Rebind(query)

	_, err := r.db.ExecContext(
		ctx,
		query,
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