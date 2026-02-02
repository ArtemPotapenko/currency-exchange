package pagination

type PageRequest struct {
	PageNumber int32
	PageSize   int32
}

type Page[T any] struct {
	Items      []T  `json:"items"`
	PageNumber int32 `json:"pageNumber"`
	PageSize   int32 `json:"pageSize"`
	Total      int   `json:"total"`
}
