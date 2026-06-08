package opensearch

type Document struct {
	Namespace string `json:"namespace"`

	ExternalID string `json:"external_id"`

	Title string `json:"title"`

	Text string `json:"text"`
}

type SearchResult struct {
	Namespace string

	ExternalID string

	Score float64
}