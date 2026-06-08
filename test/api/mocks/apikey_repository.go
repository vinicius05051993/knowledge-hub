package mocks

import (
	"context"

	"indexer/internal/apikeys"
)

type APIKeyRepository struct {
	FindByHashFunc func(
		ctx context.Context,
		hash string,
	) (*apikeys.APIKey, error)

	CreateFunc func(
		ctx context.Context,
		apiKey *apikeys.APIKey,
	) error
}

func (m *APIKeyRepository) FindByHash(
	ctx context.Context,
	hash string,
) (*apikeys.APIKey, error) {

	return m.FindByHashFunc(
		ctx,
		hash,
	)
}

func (m *APIKeyRepository) Create(
	ctx context.Context,
	apiKey *apikeys.APIKey,
) error {

	if m.CreateFunc == nil {
		return nil
	}

	return m.CreateFunc(
		ctx,
		apiKey,
	)
}