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

	defer db.Close()

	cleanupTestData(
		t,
		db,
	)

	repository :=
		apikeys.NewRepository(db)

	key := &apikeys.APIKey{
		Name:        "Test",
		Namespace:   "test",
		APIKeyHash:  uuid.NewString(),
		Permissions: "[]",
		Active:      true,
		CreatedAt:   time.Now(),
	}

	err := repository.Create(
		context.Background(),
		key,
	)

	if err != nil {
		t.Fatal(err)
	}
}