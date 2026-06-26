package documents

type SortDirection string

const (
	SortDirectionAsc  SortDirection = "asc"
	SortDirectionDesc SortDirection = "desc"
)

type OrderValueType string

const (
	OrderValueTypeString OrderValueType = "string"
	OrderValueTypeNumber OrderValueType = "number"
	OrderValueTypeDate   OrderValueType = "date"
)

type SearchRequest struct {
	Query string `json:"query"`

	Offset int `json:"offset"`

	Limit int `json:"limit"`

	Filters map[string]string `json:"filters"`

	FilterType string `json:"filter_type"`

	Order *SearchOrder `json:"order,omitempty"`
}

type SearchOrder struct {
	Field string `json:"field"`

	Direction SortDirection `json:"direction"`

	ValueType OrderValueType `json:"value_type"`
}