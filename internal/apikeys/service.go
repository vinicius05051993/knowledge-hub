package apikeys

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type RepositoryInterface interface {
	FindByHash(
		ctx context.Context,
		hash string,
	) (*APIKey, error)

	Create(
		ctx context.Context,
		apiKey *APIKey,
	) error

	DeleteByNamespace(
		ctx context.Context,
		namespace string,
	) error

	UpdateHashByNamespace(
		ctx context.Context,
		namespace string,
		hash string,
	) error
}

type Service struct {
	repository RepositoryInterface
}

func NewService(
	repository RepositoryInterface,
) *Service {

	return &Service{
		repository: repository,
	}
}

func (s *Service) Validate(
	ctx context.Context,
	apiKey string,
) (*APIKey, error) {

	hash := sha256.Sum256(
		[]byte(apiKey),
	)

	hashString := hex.EncodeToString(
		hash[:],
	)

	return s.repository.FindByHash(
		ctx,
		hashString,
	)
}

func (s *Service) Create(
	ctx context.Context,
	name string,
	namespace string,
) (string, error) {

	random := make([]byte, 32)

	_, err := rand.Read(random)

	if err != nil {
		return "", err
	}

	apiKey := "sk_live_" +
		hex.EncodeToString(random)

	hash := sha256.Sum256(
		[]byte(apiKey),
	)

	hashString := hex.EncodeToString(
		hash[:],
	)

	key := &APIKey{
		Name:        name,
		Namespace:   namespace,
		APIKeyHash:  hashString,
		Permissions: `[
			"DOCUMENT_UPSERT",
			"DOCUMENT_DELETE",
			"DOCUMENT_SEARCH"
		]`,
		Active:    true,
		CreatedAt: time.Now().UTC(),
	}

	err = s.repository.Create(
		ctx,
		key,
	)

	if err != nil {
		return "", err
	}

	return apiKey, nil
}

func (s *Service) ResetByNamespace(
	ctx context.Context,
	namespace string,
) (string, error) {

	random := make([]byte, 32)

	_, err := rand.Read(random)

	if err != nil {
		return "", err
	}

	apiKey := "sk_live_" +
		hex.EncodeToString(random)

	hash := sha256.Sum256(
		[]byte(apiKey),
	)

	hashString := hex.EncodeToString(
		hash[:],
	)

	err = s.repository.UpdateHashByNamespace(
		ctx,
		namespace,
		hashString,
	)

	if err != nil {
		return "", err
	}

	return apiKey, nil
}

func (s *Service) DeleteByNamespace(
	ctx context.Context,
	namespace string,
) error {

	return s.repository.DeleteByNamespace(
		ctx,
		namespace,
	)
}