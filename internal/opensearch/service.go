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

func (s *Service) BulkIndexDocuments(
	ctx context.Context,
	documents []*Document,
) error {

	return BulkIndexDocuments(
		ctx,
		s.client,
		documents,
	)
}

func (s *Service) Search(
	ctx context.Context,
	query string,
	offset int,
	limit int,
) ([]SearchResult, error) {

	return Search(
		ctx,
		s.client,
		query,
		offset,
		limit,
	)
}

func (s *Service) DeleteDocuments(
	ctx context.Context,
	documentKeys []string,
) error {

	return DeleteDocuments(
		ctx,
		s.client,
		documentKeys,
	)
}