package opensearch

type Document struct {
	ID string `json:"-"`

	DocumentKey string `json:"document_key"`

	Namespace string `json:"namespace"`

	ExternalID string `json:"external_id"`

	Title string `json:"title"`

	Text string `json:"text"`
}

type SearchResult struct {
	DocumentKey string

	Namespace string

	ExternalID string

	Score float64

	Highlights map[string][]string
}