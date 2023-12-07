package controller

type Page struct {
	Limit      int    `json:"limit,omitempty;query:limit"`
	Page       int    `json:"page,omitempty;query:page"`
	Sort       string `json:"sort,omitempty;query:sort"`
	TotalRows  int64  `json:"totalRows"`
	TotalPages int    `json:"totalPages"`
}

type GetPageDto struct {
	Limit  *int64  `json:"limit,omitempty;query:limit"`
	Page   *int64  `json:"page,omitempty;query:page"`
	Sort   *string `json:"sort,omitempty;query:sort"`
	Search *string `json:"search"`
}
