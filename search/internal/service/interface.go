package service

import (
	"context"
	"search/internal/models"
)

type InterfaceSearchServicee interface {
	SrvSearchGoods(ctx context.Context, req models.SearchRequest) (*models.SearchResponse, error)
}

//mockgen -source=interface.go -destination=mocks/mock_interface.go -package=mocks
