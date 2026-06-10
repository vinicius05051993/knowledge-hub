package documents

import "encoding/json"

type SearchResponse struct {
	Namespace string `json:"namespace"`

	ExternalID string `json:"external_id"`

	Title string `json:"title"`

	Text string `json:"text"`

	Payload json.RawMessage `json:"payload"`

	Highlights map[string]string `json:"highlights,omitempty"`
}