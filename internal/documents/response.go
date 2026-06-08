package documents

type UpsertDocumentsResponse struct {
	Success bool `json:"success"`

	Processed int `json:"processed"`
}