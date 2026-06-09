package documents

type UpsertDocumentsRequest struct {
	Documents []UpsertDocumentRequest `json:"documents"`
}

type UpsertDocumentRequest struct {
	ExternalID string `json:"external_id"`

	Title string `json:"title"`

	Text string `json:"text"`

	Payload map[string]interface{} `json:"payload"`
}

type DeleteDocumentsRequest struct {
	ExternalIDs []string `json:"external_ids"`
}