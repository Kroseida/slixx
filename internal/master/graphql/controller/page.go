package controller

type Page[T interface{}] struct {
	TotalRows  int64 `json:"total_rows"`
	TotalPages int   `json:"total_pages"`
	Rows       []T   `json:"rows"`
	Limit      int   `json:"limit,omitempty;query:limit"`
	Page       int   `json:"page,omitempty;query:page"`
	Sort       int   `json:"sort,omitempty;query:sort"`
}
