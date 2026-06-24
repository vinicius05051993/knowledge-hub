package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"indexer/internal/apikeys"
)

func TestAPIKeyRepositoryUpdateHashByNamespace(
	t *testing.T,
) {

	db := createDB(t)

	namespace := "test-reset"

	repository :=
		apikeys.NewRepository(db)

	t.Cleanup(func() {

		_ = repository.DeleteByNamespace(
			context.Background(),
			namespace,
		)

		_ = db.Close()
	})

	key := &apikeys.APIKey{
		Name:       "Test",
		Namespace:  namespace,
		APIKeyHash: uuid.NewString(),
		Permissions: `[
			"DOCUMENT_UPSERT",
			"DOCUMENT_DELETE",
			"DOCUMENT_SEARCH"
		]`,
		Active:    true,
		CreatedAt: time.Now(),
	}

	err := repository.Create(
		context.Background(),
		key,
	)

	if err != nil {
		t.Fatal(err)
	}

	oldHash := key.APIKeyHash
	newHash := uuid.NewString()

	err = repository.UpdateHashByNamespace(
		context.Background(),
		namespace,
		newHash,
	)

	if err != nil {
		t.Fatal(err)
	}

	apiKey, err := repository.FindByHash(
		context.Background(),
		newHash,
	)

	if err != nil {
		t.Fatal(err)
	}

	if apiKey.Namespace != namespace {

		t.Fatalf(
			"expected namespace %s got %s",
			namespace,
			apiKey.Namespace,
		)
	}

	_, err = repository.FindByHash(
		context.Background(),
		oldHash,
	)

	if err == nil {

		t.Fatal(
			"expected old hash to be invalid",
		)
	}
}

func TestAPIKeyRepositoryUpdateHashByNamespaceNotFound(
	t *testing.T,
) {

	db := createDB(t)

	repository :=
		apikeys.NewRepository(db)

	t.Cleanup(func() {
		_ = db.Close()
	})

	err := repository.UpdateHashByNamespace(
		context.Background(),
		"namespace-not-found",
		uuid.NewString(),
	)

	if err == nil {

		t.Fatal(
			"expected error for missing namespace",
		)
	}
}