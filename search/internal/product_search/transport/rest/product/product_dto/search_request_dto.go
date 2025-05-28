package product_dto

type SearchRequest struct {
	Query           string   `json:"query"`
	Categories      []string `json:"categories"`
	Brand           []string `json:"brand"`
	MinPrice        float64  `json:"minPrice"`
	MaxPrice        float64  `json:"maxPrice"`
	SortBy          string   `json:"sortBy"`
	SortOrder       string   `json:"sortOrder"`
	Page            int      `json:"page"`
	PageSize        int      `json:"pageSize"`
	HighlightFields []string `json:"highlightFields"`
}
