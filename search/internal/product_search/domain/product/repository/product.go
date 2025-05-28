package repository

import (
	entityCH "gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/domain/cluster_health/entity"
	"gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/domain/product/entity"
	entitySP "gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/domain/search_params/entity"
	entitySR "gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/domain/search_result/entity"
)

type ProductRepository interface {
	Save(product *entity.Product) error

	BulkSave(products []*entity.Product) error

	Delete(id string) error

	Search(params *entitySP.SearchParams) (*entitySR.SearchResult, error)

	Health() (*entityCH.ClusterHealth, error)

	IndicesExistsAndCreateIfMissing() error
}
