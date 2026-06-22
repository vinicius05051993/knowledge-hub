package embeddings

import "context"

type MockProvider struct {
	Embedding []float32
	Err       error
}

func (m *MockProvider) GenerateEmbedding(
	ctx context.Context,
	text string,
) ([]float32, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	if m.Embedding != nil {
		return m.Embedding, nil
	}

	return []float32{}, nil
}