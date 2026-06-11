package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"indexer/internal/apikeys"
)

func TestAPIKeyRepositoryCreate(
	t *testing.T,
) {

	db := createDB(t)

	namespace :=
		"test"

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
}