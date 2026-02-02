package dto

type CurrencyPageDto struct {
	Items      []CurrencyDto `json:"items"`
	PageNumber int32         `json:"pageNumber"`
	PageSize   int32         `json:"pageSize"`
	Total      int           `json:"total"`
}
