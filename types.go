package paginator

type Params struct {
	Page    int
	Limit   int
	OrderBy string
}

type Metadata struct {
	Page         int   `json:"page"`
	Limit        int   `json:"limit"`
	TotalRecords int64 `json:"total_records"`
	TotalPages   int   `json:"total_pages"`
	HasNext      bool  `json:"has_next"`
	HasPrev      bool  `json:"has_prev"`
}
