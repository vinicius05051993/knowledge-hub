package documents

type SearchRequest struct {
	Query string `json:"query"`

	Limit int `json:"limit"`

	Filters map[string]string `json:"filters"`
}