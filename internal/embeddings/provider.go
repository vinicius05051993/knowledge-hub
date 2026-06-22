package embeddings

import "context"

type Provider interface {
	GenerateEmbedding(
		ctx context.Context,
		text string,
	) ([]float32, error)
}