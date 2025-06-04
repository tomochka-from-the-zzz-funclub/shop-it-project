package services

import (
	"goods/internal/models"

	"github.com/google/uuid"
)

type InterfaceService interface {
	SrvCreateGoodCard(card models.GoodCard, sellerID uuid.UUID) (uuid.UUID, error)

	SrvDeleteGoodCard(cardId uuid.UUID, sellerID uuid.UUID) error

	SrvUpdateGoodCard(uuid uuid.UUID, card models.GoodCard, sellerID uuid.UUID) error

	SrvReadGood(uuid uuid.UUID, sellerID uuid.UUID) (models.Good, error)

	SrvAddCountGood(uuid uuid.UUID, number int, sellerID uuid.UUID) (int, error) //суммарное количество товара

	SrvDeleteCountGood(uuid uuid.UUID, number int, sellerID uuid.UUID) (int, error)

	SrvCreateGood(cardID uuid.UUID, quantity int, sellerID uuid.UUID) error
	SrvDeleteGood(goodID uuid.UUID, sellerID uuid.UUID) error

	SrvSearchGoods(req models.SearchRequest) (*models.SearchResponse, error)
}
