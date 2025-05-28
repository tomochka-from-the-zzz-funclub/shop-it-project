package entity

import "gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/domain/facet_bucket/entity"

type Facet struct {
	Field   string
	Buckets []entity.FacetBucket
}
