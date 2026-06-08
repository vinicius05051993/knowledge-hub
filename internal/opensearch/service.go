package opensearch

import "context"

type Service struct {
	client *Client
}

func NewService(
	client *Client,
) *Service {

	return &Service{
		client: client,
	}
}

func (s *Service) IndexDocument(
	ctx context.Context,
	document *Document,
) error {

	return IndexDocument(
		ctx,
		s.client,
		document,
	)
}