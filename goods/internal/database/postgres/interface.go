package postgresdb

import (
	"goods/internal/models"

	"github.com/google/uuid"
)

type InterfacePostgresDB interface {
	CreateGoodCard(goodCard models.GoodCard, sellerID uuid.UUID) (uuid.UUID, error)
	DeleteGoodCard(cardID uuid.UUID, sellerID uuid.UUID) error
	UpdateGoodCard(id uuid.UUID, goodCard models.GoodCard, sellerID uuid.UUID) error
	CreateGood(cardID uuid.UUID, quantity int, sellerID uuid.UUID) error
	DeleteGood(goodID uuid.UUID, sellerID uuid.UUID) error
	AddCountGood(goodID uuid.UUID, number int, sellerID uuid.UUID) (int, error)
	DeleteCountGood(goodID uuid.UUID, number int, sellerID uuid.UUID) (int, error)

	ReadGood(goodID uuid.UUID) (models.Good, error)

	SearchGoods(req models.SearchRequest) (*models.SearchResponse, error)
}
