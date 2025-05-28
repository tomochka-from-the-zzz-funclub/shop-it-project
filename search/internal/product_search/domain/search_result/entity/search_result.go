package entity

import (
	entityF "gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/domain/facet/entity"
	entityP "gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/domain/product/entity"
)

type SearchResult struct {
	Products   []*entityP.Product
	Total      int64
	Page       int
	PageSize   int
	Highlights map[string][]string
	Facets     map[string]entityF.Facet
}
