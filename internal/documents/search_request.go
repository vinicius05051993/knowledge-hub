package documents

type SearchRequest struct {
	Query string `json:"query"`

	Offset int `json:"offset"`

	Limit int `json:"limit"`

	Filters map[string]string `json:"filters"`
}