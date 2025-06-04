package service

import (
	"context"
	"search/internal/config"
	"search/internal/elastic"
	"search/internal/models"
)

type SearchService struct{}

func NewSearchService(cfg config.Config) *SearchService {
	return &SearchService{}
}

// func (s *SearchService) Search(ctx context.Context, req *models.SearchRequest) (*models.SearchResponse, error) {
// 	if req.Page < 1 {
// 		req.Page = 1
// 	}
// 	if req.PageSize <= 0 {
// 		req.PageSize = 10
// 	}

// 	products, total, err := s.db.SearchProducts(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &models.SearchResponse{
// 		Products: products,
// 		Total:    total,
// 		Page:     req.Page,
// 		PageSize: req.PageSize,
// 	}, nil
// }

func (s *SearchService) SrvSearchGoods(ctx context.Context, req models.SearchRequest) (*models.SearchResponse, error) {
	es := elastic.GetClient()
	query := elastic.BuildSearchQuery(req)

	sr, err := es.Search(
		es.Search.WithContext(ctx),
		es.Search.WithIndex("goods"),
		es.Search.WithBody(query),
	)
	if err != nil {
		return nil, err
	}
	defer sr.Body.Close()

	return elastic.ParseSearchResult(sr, req)
}
