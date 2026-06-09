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

func (s *Service) DeleteDocument(
	ctx context.Context,
	namespace string,
	externalID string,
) error {

	return DeleteDocument(
		ctx,
		s.client,
		namespace,
		externalID,
	)
}