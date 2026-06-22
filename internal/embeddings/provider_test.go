package embeddings

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestMockProvider_ReturnsConfiguredEmbedding(t *testing.T) {
	t.Parallel()

	expected := []float32{
		0.1,
		0.2,
		0.3,
	}

	provider := &MockProvider{
		Embedding: expected,
	}

	result, err := provider.GenerateEmbedding(
		context.Background(),
		"hello world",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf(
			"expected %v, got %v",
			expected,
			result,
		)
	}
}

func TestMockProvider_ReturnsConfiguredError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("embedding provider failed")

	provider := &MockProvider{
		Err: expectedErr,
	}

	result, err := provider.GenerateEmbedding(
		context.Background(),
		"hello world",
	)

	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(err, expectedErr) {
		t.Fatalf(
			"expected error %v, got %v",
			expectedErr,
			err,
		)
	}

	if result != nil {
		t.Fatalf(
			"expected nil embedding, got %v",
			result,
		)
	}
}

func TestMockProvider_ReturnsEmptyEmbeddingByDefault(t *testing.T) {
	t.Parallel()

	provider := &MockProvider{}

	result, err := provider.GenerateEmbedding(
		context.Background(),
		"hello world",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil embedding")
	}

	if len(result) != 0 {
		t.Fatalf(
			"expected empty embedding, got %v",
			result,
		)
	}
}