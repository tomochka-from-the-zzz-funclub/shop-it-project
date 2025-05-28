package use_cases

import (
	"fmt"
	"log"

	"gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/transport/rest/product/product_dto"

	entityCH "gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/domain/cluster_health/entity"
	entityP "gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/domain/product/entity"
	"gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/domain/product/repository"
	entitySP "gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/domain/search_params/entity"
	entitySR "gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/domain/search_result/entity"
)

type ProductService struct {
	Repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	log.Println("[ProductService] Initializing...")

	if err := repo.IndicesExistsAndCreateIfMissing(); err != nil {
		log.Printf("[ProductService][ERROR] Failed to ensure index: %v", err)
	}

	log.Println("[ProductService] Initialized")
	return &ProductService{Repo: repo}
}

func (s *ProductService) IndexProduct(p *entityP.Product) error {
	log.Printf("[ProductService] IndexProduct called for ID=%s", p.ID)
	if p == nil {
		err := fmt.Errorf("product is nil")
		log.Printf("[ProductService][ERROR] %v", err)
		return err
	}
	if err := s.Repo.Save(p); err != nil {
		log.Printf("[ProductService][ERROR] Save failed for ID=%s: %v", p.ID, err)
		return err
	}
	log.Printf("[ProductService] Product indexed successfully ID=%s", p.ID)
	return nil
}

func (s *ProductService) BulkIndexProducts(products []*entityP.Product) error {
	log.Printf("[ProductService] BulkIndexProducts called for %d products", len(products))
	if len(products) == 0 {
		log.Println("[ProductService] No products to index")
		return nil
	}
	if err := s.Repo.BulkSave(products); err != nil {
		log.Printf("[ProductService][ERROR] BulkSave failed: %v", err)
		return err
	}
	log.Printf("[ProductService] Bulk indexing succeeded for %d products", len(products))
	return nil
}

func (s *ProductService) DeleteProduct(id string) error {
	log.Printf("[ProductService] DeleteProduct called for ID=%s", id)
	if id == "" {
		err := fmt.Errorf("product ID is empty")
		log.Printf("[ProductService][ERROR] %v", err)
		return err
	}
	if err := s.Repo.Delete(id); err != nil {
		log.Printf("[ProductService][ERROR] Delete failed for ID=%s: %v", id, err)
		return err
	}
	log.Printf("[ProductService] Product deleted successfully ID=%s", id)
	return nil
}

func (s *ProductService) SearchProducts(req *product_dto.SearchRequest) (*entitySR.SearchResult, error) {
	log.Printf("[ProductService] SearchProducts called with params: %+v", req)
	if req == nil {
		err := fmt.Errorf("request params are nil")
		log.Printf("[ProductService][ERROR] %v", err)
		return nil, err
	}

	params := &entitySP.SearchParams{
		Query:           req.Query,
		Categories:      req.Categories,
		Brand:           req.Brand,
		MinPrice:        req.MinPrice,
		MaxPrice:        req.MaxPrice,
		SortBy:          req.SortBy,
		SortOrder:       req.SortOrder,
		Page:            req.Page,
		PageSize:        req.PageSize,
		HighlightFields: req.HighlightFields,
	}

	if params.PageSize == 0 {
		params.PageSize = 10
	}
	if params.Page == 0 {
		params.Page = 1
	}

	result, err := s.Repo.Search(params)
	if err != nil {
		log.Printf("[ProductService][ERROR] Search failed: %v", err)
		return nil, err
	}

	log.Printf("[ProductService] SearchProducts succeeded: found %d items", result.Total)
	return result, nil
}

func (s *ProductService) GetClusterHealth() (*entityCH.ClusterHealth, error) {
	log.Println("[ProductService] GetClusterHealth called")
	health, err := s.Repo.Health()
	if err != nil {
		log.Printf("[ProductService][ERROR] Health check failed: %v", err)
		return nil, err
	}
	log.Printf("[ProductService] Cluster health: status=%s", health.Status)
	return health, nil
}
